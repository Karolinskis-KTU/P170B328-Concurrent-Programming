#include "cuda_runtime.h"
#include "device_launch_parameters.h"

#include <fstream>
#include <iostream>
#include <vector>
#include <stdio.h>

#include "json/single_include/nlohmann/json.hpp"

using json = nlohmann::json;

struct Car {
    char name[50];
    int fuel_tank_size;
    float fuel_efficiency;
};

struct Result {
    char data[60];
};

std::vector<Car> readFile(const std::string& filename) {
    std::ifstream file(filename);
    if (!file.is_open()) {
		std::cout << "Could not open file: " << filename << std::endl;
        exit(EXIT_FAILURE);
	}

    nlohmann::json jsonData;
    file >> jsonData;
    file.close();

    std::vector<Car> cars;
    for (const auto& carData : jsonData["cars"]) {
        Car car;
        strcpy(car.name, carData["name"].get<std::string>().c_str());
        car.fuel_tank_size = carData["fuel_tank_size"].get<int>();
        car.fuel_efficiency = carData["fuel_efficiency"].get<float>();
        cars.push_back(car);
    }

    return cars;
}

void printResults(const Result* results, int size) {
    int counter = 0;
    for (int i = 0; i < size; i++) {
        if (results[i].data[0] != '\0') {
            std::cout << results[i].data << std::endl;
            counter++;
        }
    }
    std::cout << "Number of results: " << counter << std::endl;
}

void writeResultsToFile(const Result* results, int size, std::string filename) {
	std::ofstream file(filename);
	if (!file.is_open()) {
		std::cout << "Could not open file: " << filename << std::endl;
		exit(EXIT_FAILURE);
	}

	int counter = 0;
	for (int i = 0; i < size; i++) {
		if (results[i].data[0] != '\0') {
			file << results[i].data << std::endl;
			counter++;
		}
	}
    file << "Number of results: " << counter;
	file.close();

    std::cout << "Successfully wrote results to file: " << filename << std::endl;
}

void writeResultsToConsole(const Result* results, int size) {
    char ans;
    
    std::cout << "Do you want to see the results in the console? (Enter 'Y' for Yes, 'N' for No): ";
    std::cin >> ans;
    ans = std::toupper(ans);

    if (ans != 'Y') {
        return;
    }

    int counter = 0;
    for (int i = 0; i < size; i++) {
        if (results[i].data[0] != '\0') {
            std::cout << results[i].data << std::endl;
            counter++;
        }
    }

    std::cout << "Number of results: " << counter << std::endl;
}

void checkCudaDevice() {
    int deviceCount;
    cudaGetDeviceCount(&deviceCount);

    if (deviceCount == 0) {
		std::cout << "There is no CUDA device" << std::endl;
		exit(EXIT_FAILURE);
	}

    for (int i = 0; i < deviceCount; ++i) {
        cudaDeviceProp deviceProp;
        cudaGetDeviceProperties(&deviceProp, i);
        std::cout << "Device " << i << ": " << deviceProp.name << std::endl;
    }
}

__device__ char convertToRating(float efficiency) {
    if (efficiency > 30.0f) {
        return 'A';
    }
    else if (efficiency >= 25.0f) {
        return 'B';
    } 
    else if (efficiency >= 20.0f) {
		return 'C';
	}
	else if (efficiency >= 15.0f) {
		return 'D';
	}
	else {
		return 'E';
	}
}

__global__ void filterAndSortCars(Car* cars, int size, int tankSizeThreshold, Result* results, int resultsSize) {
    int index = blockIdx.x * blockDim.x + threadIdx.x;

    if (index < size) {
        if (cars[index].fuel_tank_size > tankSizeThreshold) {

            Result result;

            // Copy name to result.data
            const char* namePtr = cars[index].name;
            char* resultPtr = result.data;
            while (*namePtr != '\0') {
				*resultPtr = *namePtr;
				++namePtr;
				++resultPtr;
			}

            *resultPtr = '-';
            ++resultPtr;

            // Convert fuel_efficiency to rating and append to result.data
            *resultPtr = convertToRating(cars[index].fuel_efficiency);
            ++resultPtr;

            // Append fuel tank size to result.data
            int tankSize = cars[index].fuel_tank_size;
            if (tankSize >= 100) {
                *resultPtr = '1';
                ++resultPtr;
                *resultPtr = '0' + (tankSize % 100) / 10;
                ++resultPtr;
                *resultPtr = '0' + tankSize % 10;  // Increment resultPtr for the last digit
                ++resultPtr;
            }
            else {
                *resultPtr = '0' + tankSize / 10;
                ++resultPtr;
                *resultPtr = '0' + tankSize % 10;  // Increment resultPtr for the last digit
                ++resultPtr;
            }

            // Find the first free slot in the results array
            int resultIndex = index % resultsSize;

            while (atomicCAS((int*)&results[resultIndex].data[0], 0, 1) != 0) {
                resultIndex = (resultIndex + 1) % resultsSize;
            }

            results[resultIndex] = result;
        }
    }
}

int main() {
    // Check and print CUDA device information
    checkCudaDevice();

    std::string outputFile = "IFF-1-1_PaulaviciusK_L3_res.txt";

    // User input to select file to read
    std::string inputFile;
    std::cout << "Select file (1 to 3): ";
    std::cin >> inputFile;
    if (inputFile == "1") {
        inputFile = "IFF-1-1_PaulaviciusK_L3_dat_1.json";
    }
    else if (inputFile == "2") {
        inputFile = "IFF-1-1_PaulaviciusK_L3_dat_2.json";
    }
    else if (inputFile == "3") {
        inputFile = "IFF-1-1_PaulaviciusK_L3_dat_3.json";
    }
    else {
        std::cout << "Invalid input" << std::endl;
        exit(EXIT_FAILURE);
    }

    std::cout << "Reading file: " << inputFile << std::endl;

    // Host variables
    std::vector<Car> hostCars = readFile(inputFile);
    const int dataSize = hostCars.size();
    const int fuelTankSizeThreshold = 60;

    // Device variables
    Car* deviceCars;
    Result* hostResults = new Result[dataSize];
    Result* deviceResults;

    // Ensure at least two blocks and the number of threads per block is a multiple of 32
    const int threadsPerBlock = 32 * 2;
    const int blocksPerGrid = (dataSize + threadsPerBlock - 1) / threadsPerBlock;

    // Allocate memory on GPU
    cudaMalloc((void**)&deviceCars, dataSize * sizeof(Car));
    cudaMalloc((void**)&deviceResults, dataSize * sizeof(Result));

    // Copy data from host to device
    cudaMemcpy(deviceCars, hostCars.data(), dataSize * sizeof(Car), cudaMemcpyHostToDevice);

    // Initialize results array on host
    memset(hostResults, 0, dataSize * sizeof(Result));

    // Copy results array from host to device
    cudaMemcpy(deviceResults, hostResults, dataSize * sizeof(Result), cudaMemcpyHostToDevice);

    // Launch CUDA kernel
    filterAndSortCars<<<blocksPerGrid, threadsPerBlock>>>(deviceCars, dataSize, fuelTankSizeThreshold, deviceResults, dataSize);
    cudaDeviceSynchronize(); // Add synchronization to ensure the kernel is finished before copying data back

    // Copy data from device to host
    cudaMemcpy(hostResults, deviceResults, dataSize * sizeof(Result), cudaMemcpyDeviceToHost);

    // Print results
    writeResultsToFile(hostResults, dataSize, outputFile);

    // Ask to write to console
    writeResultsToConsole(hostResults, dataSize);

    // Free memory on GPU and host
    cudaFree(deviceCars);
    cudaFree(deviceResults);

    return 0;
}
