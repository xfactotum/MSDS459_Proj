import requests
from bs4 import BeautifulSoup
import json

def scrape_akc_data(group_names):
    all_data = []

    for group_name in group_names:
        # Fetch description
        response = requests.get(f"https://www.akc.org/dog-breeds/{group_name}/")
        if response.status_code == 200:
            soup = BeautifulSoup(response.content, "html.parser")
            description = soup.find("div", class_="read-more__long").p.text.strip()
        else:
            print(f"Failed to retrieve the web page ({group_name}). Status code: {response.status_code}")
            description = "Description not available."

        # Fetch dog breeds
        dog_breeds = []
        for page_num in range(1, 4):
            url = f"https://www.akc.org/dog-breeds/{group_name}/page/{page_num}/" if page_num > 1 else f"https://www.akc.org/dog-breeds/{group_name}/"
            response = requests.get(url)

            if response.status_code == 200:
                soup = BeautifulSoup(response.content, "html.parser")
                breed_tiles = soup.find_all("h3", class_="breed-type-card__title")
                if not breed_tiles:
                    break  # No more breeds on the page

                dog_breeds.extend([breed.text.strip() for breed in breed_tiles])
            elif response.status_code == 404:
                print(f"Page not found: {url}")
            else:
                print(f"Failed to retrieve the web page ({url}). Status code: {response.status_code}")
                break

        # Create a dictionary for this group
        group_data = {
            "page_header_title": f"{group_name.capitalize()} Group",
            "description": description,
            "dog_breeds": dog_breeds
        }

        all_data.append(group_data)

    # Save the combined data to a JSON file
    with open("akc_breed_groups.json", "w") as json_file:
        json.dump(all_data, json_file, indent=4)
        
# List of various group names on akc website:
group_names = ['sporting', 'hound', 'working', 'terrier', 'toy', 'non-sporting', 'herding', 
               'miscellaneous-class', 'smallest-dog-breeds', 'medium-dog-breeds', 'largest-dog-breeds',
               'smartest-dogs', 'hypoallergenic-dogs', 'best-family-dogs', 'best-guard-dogs', 
               'best-dogs-for-kids', 'best-dogs-for-apartment-dwellers', 'hairless-dog-breeds']

# Scrape data for all groups and save to JSON
scrape_akc_data(group_names)