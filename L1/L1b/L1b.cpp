#include <iostream>
#include <iomanip> 
#include <string>
#include <fstream>
#include <omp.h>
#include <cstdint>
#include <functional>
#include <sstream>
#include <thread>
#include <vector>
#include "nlohmann/json.hpp"

using json = nlohmann::json;

// Define a Car class
class Car {
public:
    std::string name;
    int fuelTankSize;
    double fuelEfficiency;

    /// <summary>
    /// Parse JSON data and populate the object
    /// </summary>
    /// <param name="dataString">Data to parse to Car object</param>
    void fromJson(json dataString) {
        std::string temp = dataString["name"];
        name = temp;
        fuelTankSize = dataString["fuel_tank_size"];
        fuelEfficiency = dataString["fuel_efficiency"];
    }

    /// <summary>
    /// Generate a hashcode for the Car object
    /// </summary>
    /// <returns>An int of the car object hash code</returns>
    int hashCode() const {
        std::string data = name + std::to_string(fuelTankSize) + std::to_string(static_cast<int>(fuelEfficiency));
        std::hash<std::string> hash_fn;
        return static_cast<int>(hash_fn(data));
    }
};

/// <summary>
/// Define a Monitor class to manage Car objects
/// </summary>
class Monitor {
private:
    std::vector<Car> list;
    omp_lock_t lock;

public:
    int iSum = 0;
    double dSum = 0;

    Monitor() {
        omp_init_lock(&lock);
    }

    ~Monitor() {
        omp_destroy_lock(&lock);
    }

    /// <summary>
    /// Add a Car object to the list in a thread-safe manner
    /// </summary>
    /// <param name="car">Car to add</param>
    void Add(Car& car) {
        #pragma omp critical
        {
            auto it = std::lower_bound(list.begin(), list.end(), car, [](const Car& a, const Car& b) {
                return a.name > b.name;
                });
            list.insert(it, car);
        }
    }

    /// <summary>
    /// Remove and return a Car object from the list in a thread-safe manner
    /// </summary>
    /// <returns>Returns the removed Car object</returns>
    Car Pop() {
        Car car;

        #pragma omp critical
        {
            if (!list.empty()) {
                car = list.back();
                list.pop_back();
            }
        }
        return car;
    }

    /// <summary>
    /// Gets the count of Car objects in the list
    /// </summary>
    /// <returns>Returns the ammount of Car objects in the list</returns>
    int Count() const {
        return list.size();
    };
};


/// <summary>
/// Utility class for file I/O
/// </summary>
class IO {
public:
    /// <summary>
    /// Read Car data from a JSON file
    /// </summary>
    /// <param name="fileName">file path</param>
    /// <returns>Returns a list of all cars read from the JSON file</returns>
    static std::vector<Car> ReadFile(std::string& fileName) {
        std::vector<Car> cars;
        std::ifstream stream;
        stream.open(fileName);

        json allCarsJson = json::parse(stream);
        auto allCars = allCarsJson["cars"];
        for (const json& car : allCars) {
            Car tempCar;
            tempCar.fromJson(car);
            cars.push_back(tempCar);
        }
        stream.close();
        return cars;
    }

    /// <summary>
    /// Print Car data and result to a file
    /// </summary>
    /// <param name="fileName">File to write to</param>
    /// <param name="list">List of all cars</param>
    /// <param name="iSum">Sum of all integers</param>
    /// <param name="dSum">Sum of all doubles</param>
    static void printResult(std::string& fileName, Monitor& list, int iSum, double dSum) {
        const char* filePath = fileName.c_str();

        std::ofstream file;
        file.open(filePath, std::ios_base::app);
        file << std::string(30, ' ') <<  "Results" << std::endl;
        file << "+"<< std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << "|" << std::setw(5) << "# |" << std::setw(21) << "Name |" << std::setw(21) << "Fuel Efficiency |" << std::setw(22) << "Fuel tank size |" << std::endl;
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        int i = 0;
        while (list.Count() > 0) {
            Car temp = list.Pop();
            file << "|" << std::setw(3) << i << " |" << std::setw(20) << temp.name << "|" << std::setw(20) << temp.fuelEfficiency << "|" << std::setw(20) << temp.fuelTankSize << " |" << std::endl;
            i++;
        }
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << "Total sum of int fields:" << iSum << std::endl;
        file << "Total sum of double fields:" << dSum << std::endl;
    }

    /// <summary>
    /// Print Car data to a file
    /// </summary>
    /// <param name="fileName">File to write to</param>
    /// <param name="cars">List of all files</param>
    static void printResult(std::string& fileName, std::vector<Car>& cars) {
        const char* filePath = fileName.c_str();

        remove(filePath);

        std::ofstream file;
        file.open(filePath);
        file << std::string(30, ' ') << "Data" << std::endl;
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << "|" << std::setw(5) << "# |" << std::setw(21) << "Name |" << std::setw(21) << "Fuel Efficiency |" << std::setw(22) << "Fuel tank size |" << std::endl;
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        for (int i = 0; i < cars.size(); i++)
        {
            Car temp = cars[i];
            file << "|" << std::setw(3) << i << " |" << std::setw(20) << temp.name << "|" << std::setw(20) << temp.fuelEfficiency << "|" << std::setw(20) << temp.fuelTankSize << " |" << std::endl;
        }
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << std::endl;
    }
};

/// <summary>
/// Utility class for mathematical operations
/// </summary>
class Utils {
public:
    /// <summary>
    /// Find the closest Fibonacci number to a given integer
    /// </summary>
    /// <param name="n">Integer to find the closest Fibonacci number to</param>
    /// <returns>Returns the closest Fibonacci number</returns>
    static int closestFibonacci(int n) {
        int a = 0;
        int b = 1;
        int fib = 0;

        while (fib <= n) {
            fib = a + b;
            a = b;
            b = fib;
        }

        if (n < fib - n) {
            return a;
        }
        else {
            return b;
        }
    }
};


/// <summary>
/// Function to process Car objects and populate the Monitor
/// </summary>
/// <param name="cars">List of all Car objects</param>
/// <param name="monitor">Monitor to populate the data to</param>
/// <param name="iSum">Sum of all integers</param>
/// <param name="dSum">Sum of all doubles</param>
void execute(std::vector<Car>& cars, Monitor& monitor, int& iSum, double& dSum) {
    for (Car& car : cars) {
        int hashCode = car.hashCode();
        int fib = Utils::closestFibonacci(hashCode);

        int carSum = 0;
        while (fib != 0) {
            carSum += fib % 10;
            fib /= 10;
        }

        if (carSum % 2 == 0) {
            monitor.Add(car);
            iSum += car.fuelTankSize;
            dSum += car.fuelEfficiency;
        }
    }
}

int main() {
    std::string inputFile = "";
    std::string outputFile = "IFF-1-1_PaulaviciusK_L1_res.txt";

    int inputNum;
    std::cout << "Choose which file you want to use: ";
    std::cin >> inputNum;

    switch (inputNum) {
        case 1:
            inputFile = "IFF-1-1_PaulaviciusK_L1_dat_1.json";
            break;
        case 2:
            inputFile = "IFF-1-1_PaulaviciusK_L1_dat_2.json";
            break;
        case 3:
            inputFile = "IFF-1-1_PaulaviciusK_L1_dat_3.json";
            break;
            
        default:
            std::cout << "Unknown number provided: " << inputNum << std::endl;
            return 1;
    }

    std::cout << "Input file: " << inputFile << std::endl;


    auto cars = IO::ReadFile(inputFile);
    int numCars = cars.size();
    const int numThreads = std::max(2, static_cast<int>(std::ceil(static_cast<double>(numCars) / 2)));
    std::vector<std::vector<Car>> threadData(numThreads);
    
    int carsInThread = numCars / numThreads;
    int remainingCars = numCars % numThreads;

    int start = 0;

    for (int i = 0; i < numThreads; i++) {
        int end = start + carsInThread;

        if (remainingCars > 0) {
            end++;
            remainingCars--;
        }

        std::vector<Car> batch;
        for (int j = start; j < end; j++) {
            batch.push_back(cars[j]);
        }

        threadData[i] = std::move(batch);

        start = end;
    }

    Monitor monitor;
    int iSum = 0;
    double dSum = 0.00;

    #pragma omp parallel num_threads(numThreads) reduction(+:iSum) reduction(+:dSum)
    {
        int threadID = omp_get_thread_num();
        execute(threadData[threadID], monitor, iSum, dSum);

    }

    // Print results to the output file
    IO::printResult(outputFile, cars);
    IO::printResult(outputFile, monitor, iSum, dSum);


    return 0;
}