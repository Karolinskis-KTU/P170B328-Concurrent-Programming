import numpy as np
import json
import time
from multiprocessing import Pool
import matplotlib.pyplot as plt

# Settings
max_processes = 8
learning_rate = 0.01
num_iterations = 150

# Globals
runtime = []
CITY_SIZE_X = (0, 0)
CITY_SIZE_Y = (0, 0)
existing_stores_count = 0
new_stores_count = 0


# Read data from file
def read_data(filename):
    global CITY_SIZE_X, CITY_SIZE_Y, existing_stores, new_stores, existing_stores_count, new_stores_count

    print(f"Reading data from {filename}")
    with open("Data/"+filename, 'r') as f:
        data = json.load(f)
    print(f"Data read successfully.")
    print(f"City size: {data['city_size_x']}, {data['city_size_y']}")
    print(f"Existing stores: {len(data['existing_stores'])}")
    print(f"New stores: {len(data['new_stores'])}")

    CITY_SIZE_X = data["city_size_x"]
    CITY_SIZE_Y = data["city_size_y"]
    existing_stores = np.array(data["existing_stores"])
    new_stores = np.array(data["new_stores"])
    existing_stores_count = len(existing_stores)
    new_stores_count = len(new_stores)


def distance_to_location_cost(x1, y1, x2, y2):
    return np.exp(-0.2 * ((x1 - x2)**2 + (y1 - y2)**2))

def softplus(x):
    return np.log1p(np.exp(-np.abs(x))) + np.maximum(x, 0)

def distance_to_boundary_cost(x, y):
    dist_x = min(abs(x - CITY_SIZE_X[0]), abs(x - CITY_SIZE_X[1]))
    dist_y = min(abs(y - CITY_SIZE_Y[0]), abs(y - CITY_SIZE_Y[1]))
    dist = min(dist_x, dist_y)
    if dist == 0:
        return np.inf  # return a high cost if the store is on the city boundary
    else:
        return 1 / softplus(0.25 * dist**2)  # return a higher cost when the store is closer to the city boundary
    
def closest_point_to_edge(x, y):
    x = max(CITY_SIZE_X[0], min(x, CITY_SIZE_X[1]))
    y = max(CITY_SIZE_Y[0], min(y, CITY_SIZE_Y[1]))
    
    if y < (CITY_SIZE_Y[0] + CITY_SIZE_Y[1]) / 2:
        return (x, CITY_SIZE_Y[0])
    else:
        return (x, CITY_SIZE_Y[1]) if y >= (CITY_SIZE_Y[0] + CITY_SIZE_Y[1]) / 2 else (CITY_SIZE_X[0], y) if x < (CITY_SIZE_X[0] + CITY_SIZE_X[1]) / 2 else (CITY_SIZE_X[1], y)

def objective(existing_locations, new_locations):
    total_cost = 0
    for new_location in new_locations:
        distance_to_existing_stores = np.mean([distance_to_location_cost(new_location[0], new_location[1], existing_location[0], existing_location[1]) for existing_location in existing_locations])
        distance_to_boundary = distance_to_boundary_cost(new_location[0], new_location[1])
        total_cost += distance_to_existing_stores + distance_to_boundary
    return total_cost

def calculate_gradient(args):
    existing_stores, new_stores, j, k = args
    new_stores_count = new_stores.shape[0]  # get the number of new stores
    perturbation = np.zeros((new_stores_count,2))
    perturbation[j,k] = 1e-5
    f_plus = objective(existing_stores, new_stores + perturbation)
    f_minus = objective(existing_stores, new_stores - perturbation)
    return (f_plus - f_minus) / (2 * 1e-5)

def visualisation(existing_stores, new_stores, title):
    # Visualize the results
    plt.figure(figsize=(10, 10))
    plt.scatter(existing_stores[:, 0], existing_stores[:, 1], color='blue', marker='o', label=f'Miesto parduotuves: {len(existing_stores)}')
    plt.scatter(new_stores[:, 0], new_stores[:, 1], color='red', marker='x', label=f'Naujos parduotuves: {len(new_stores)}')
    # Draw city boundaries
    plt.plot([CITY_SIZE_X[0], CITY_SIZE_X[1]], [CITY_SIZE_Y[0], CITY_SIZE_Y[0]], color='black')
    plt.plot([CITY_SIZE_X[0], CITY_SIZE_X[1]], [CITY_SIZE_Y[1], CITY_SIZE_Y[1]], color='black')
    plt.plot([CITY_SIZE_X[0], CITY_SIZE_X[0]], [CITY_SIZE_Y[0], CITY_SIZE_Y[1]], color='black')
    plt.plot([CITY_SIZE_X[1], CITY_SIZE_X[1]], [CITY_SIZE_Y[0], CITY_SIZE_Y[1]], color='black')
    plt.xlabel('X')
    plt.ylabel('Y')
    plt.legend(loc='upper right')
    plt.title(f'{title}')
    plt.grid(True)
    plt.show()

def plot_objective_function(num_iterations, objective_values):
    """
    Create and plot the objective function graph.

    Parameters:
    num_iterations (int): The number of iterations.
    objective_values (list): The values of the objective function at each iteration.
    """
    plt.figure()
    plt.plot(range(num_iterations), objective_values, label="Tikslo funkcija")
    plt.xlabel('Iteracijos')
    plt.ylabel('Tikslo funkcija')
    plt.legend(loc='upper right')
    plt.title('Tikslo funkcijos priklausomybė nuo iteracijos')
    plt.grid(True)
    plt.show()

def plot_runtime(runtimes, dots, iterations):
    plt.figure()
    plt.plot(range(1, len(runtimes) + 1), runtimes, marker='o')
    plt.xlabel('Procesų skaičius')
    plt.ylabel('Laikas, s')
    plt.suptitle('Laiko priklausomybė nuo procesų skaičiaus')
    plt.title('Taškai: ' + str(dots) + ', Iteracijos: ' + str(iterations))
    plt.grid(True)
    plt.show()    

if __name__ == "__main__":
    # Read data from file
    read_data("existing18_new25.json")


    # visualisation(existing_stores, new_stores, "Pradinės parduotuvių vietos")
    print(objective(existing_stores, new_stores))

    for processes in range(1, max_processes + 1):
        start_time = time.time()
        with Pool(processes=processes) as pool:
            objective_values = []
            print(f"Starting calculations for pool with {processes} processes...")
            for i in range(num_iterations):
                print(f"\rStarting iteration {i+1}...", end='', flush=True)
                args = [(existing_stores, new_stores, j, k) for j in range(new_stores_count) for k in range(2)]
                grad = np.array(pool.map(calculate_gradient, args)).reshape(new_stores_count, 2)
                new_stores -= learning_rate * grad
                objective_value = objective(existing_stores, new_stores)
                objective_values.append(objective_value)
        end_time = time.time()
        runtime.append(end_time - start_time)
        print(f"\nFinished pool with {processes} processes in {runtime[-1]} seconds")


    print(objective(existing_stores, new_stores))
    #visualisation(existing_stores, new_stores, "Galutinės parduotuvių vietos")

    print(f"Esamos parduotuvės: {existing_stores_count}")
    print(existing_stores)
    print(f"Naujos parduotuvės: {new_stores_count}")
    print(new_stores)

    plot_runtime(runtime, new_stores_count, num_iterations)
    # plot_objective_function(num_iterations, objective_values)