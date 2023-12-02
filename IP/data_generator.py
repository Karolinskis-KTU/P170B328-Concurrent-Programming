import random
import numpy as np
import json

CITY_SIZE_X = (-10, 10)
CITY_SIZE_Y = (-10, 10)
existing_stores_count = np.random.randint(3, 20)
new_stores_count = np.random.randint(3, 100)
new_stores_count = 500

if __name__ == "__main__":
    print(f"Generating {existing_stores_count} existing stores...")
    existing_stores = np.random.rand(existing_stores_count, 2) * (CITY_SIZE_X[1] - CITY_SIZE_X[0]) + CITY_SIZE_X[0]
    print(f"Generating {new_stores_count} new stores...")
    new_stores = np.random.rand(new_stores_count, 2) * (CITY_SIZE_X[1] - CITY_SIZE_X[0]) + CITY_SIZE_X[0]

    existing_stores_list = existing_stores.tolist()
    new_stores_list = new_stores.tolist()

    # Convert the lists to a dictionary
    data = {
        "city_size_x": CITY_SIZE_X,
        "city_size_y": CITY_SIZE_Y,
        "existing_stores": existing_stores_list,
        "new_stores": new_stores_list
    }

    filename = f"existing{existing_stores_count}_new{new_stores_count}.json"

    print(f"Saving as {filename}")

    with open("Data/"+filename, 'w') as f:
        json.dump(data, f, indent=4)