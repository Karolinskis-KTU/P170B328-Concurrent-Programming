#include <iostream>
#include <iomanip> 
#include <string>
#include <fstream>
#include <omp.h>
#include "nlohmann/json.hpp"
using namespace std;
using json = nlohmann::json;

class Car {
public:
    string name;
    int fuelTankSize;
    double fuelEfficiency;
    void fromJson(json dataString) {
        std::string temp = dataString["name"];
        name = temp;
        fuelTankSize = dataString["fuel_tank_size"];
        fuelEfficiency = dataString["fuel_efficiency"];
    }
    int hashCode() const {
        std::string data = name + std::to_string(fuelTankSize) + std::to_string(static_cast<int>(fuelEfficiency));
        std::hash<std::string> hash_fn;
        return static_cast<int>(hash_fn(data));
    }
};

class Monitor {
public:
    int count;
    int capacity;
    Car* list;

    Monitor(int arraySize) {
        list = new Car[arraySize];
        count = 0;
        capacity = arraySize;
    }

    void add(Car& car) {
        bool carAdded = false;

        int threadId = omp_get_thread_num();

        #pragma omp critical(monitor_add)
        {
            for (int i = 0; i < count; i++) {
                if (list[i].name == car.name &&
                    list[i].fuelTankSize == car.fuelTankSize &&
                    list[i].fuelEfficiency == car.fuelEfficiency) {
                    carAdded = true;  // Car is already in the list
                    break;
                }
            }

            if (!carAdded && count != capacity) {
                int i;
                for (i = count - 1; (i >= 0 && list[i].name < car.name); i--) {
                    list[i + 1] = list[i];
                }
                list[i + 1] = car;
                count++;
                carAdded = true;
            }
        }

        if (!carAdded) {
            cout << "MonitorConsumer is full" << endl;
            exit(1);
        }
    }


    Car Get(int index) {
        Car temp;

        #pragma omp critical
        {
            temp = list[index];

        }
        return temp;
    }
};

vector<Car> ReadFile(string fileName) {
    vector<Car> cars;
    ifstream stream;
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

void printResult(string fileName, Monitor& list, double sumResult) {
    const char* filePath = fileName.c_str();
    
    ofstream file;
    file.open(filePath, std::ios_base::app);
    file << "Results" << endl;
    file << "----------------------------------------------------------------------" << endl;
    file << setw(6) << "#|" << setw(21) << "Name |" << setw(21) << "Fuel Efficiency |" << setw(22) << "Fuel tank size |" << endl;
    for (int i = 0; i < list.count; i++) {
        Car temp = list.Get(i);
        file << setw(5) << i << "|" << setw(20) << temp.name << "|" << setw(20) << temp.fuelEfficiency << "|" << setw(20) << temp.fuelTankSize << " |" << endl;
    }
    file << "----------------------------------------------------------------------" << endl;
    file << "Total sum of int and float fields:" << sumResult << endl;
    file << endl;
}

void printResult(string fileName, vector<Car> cars) {
    const char* filePath = fileName.c_str();

    remove(filePath);

    ofstream file;
    file.open(filePath);
    file << "Data" << endl;
    file << "----------------------------------------------------------------------" << endl;
    file << setw(6) << "#|" << setw(21) << "Name |" << setw(21) << "Fuel Efficiency |" << setw(22) << "Fuel tank size |" << endl;
    for (int i = 0; i < cars.size(); i++)
    {
        Car temp = cars[i];
        file << setw(5) << i << "|" << setw(20) << temp.name << "|" << setw(20) << temp.fuelEfficiency << "|" << setw(20) << temp.fuelTankSize << " |" << endl;
    }
    file << "----------------------------------------------------------------------" << endl;
    file << endl;
}

int closestFibonacci(int n) {
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
    } else {
        return b;
    }
}

int main() {
    string inputFile = "";
    string outputFile = "IFF-1-1_PaulaviciusK_L1_res.txt";

    int inputNum;
    std::cout << "Choose which file you want to use: ";
    cin >> inputNum;

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
            std::cout << "Unknown number provided: " << inputNum << endl;
            return 1;
    }

    std::cout << "Input file: " << inputFile << endl;


    auto cars = ReadFile(inputFile);
    Monitor monitor(cars.size() + 1);
    double sum = 0;
    #pragma omp parallel shared(cars, monitor, sum) num_threads(4)
    {
        auto total_threads = omp_get_num_threads();
        auto chunk_size = cars.size() / total_threads;
        auto thread_number = omp_get_thread_num();
        auto start_index = chunk_size * thread_number;
        auto end_index = thread_number == total_threads - 1 ? cars.size() : ((thread_number + 1) * chunk_size);
        for (int i = 0; i < end_index; i++)
        {
            Car tempCar = cars[i];
            int hashCode = tempCar.hashCode();
            int fib = closestFibonacci(hashCode);
            
            int carSum = 0;
            while (fib != 0) {
                carSum += fib % 10;
                fib /= 10;
            }

            // If the sum is even, add it
            if (carSum % 2 == 0) {
                monitor.add(tempCar);
            }
            #pragma omp critical
            {
                sum += tempCar.fuelEfficiency + tempCar.fuelTankSize;
            }
        }
    }
    
    printResult(outputFile, cars);
    printResult(outputFile, monitor, sum);


    return 0;
}