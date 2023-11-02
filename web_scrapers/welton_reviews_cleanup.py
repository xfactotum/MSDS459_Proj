import pandas as pd
from fuzzywuzzy import fuzz
from fuzzywuzzy import process

# Load the JSON and CSV data into dataframes
json_df = pd.read_json('web_scrapers/welton_reviews.json')
csv_df = pd.read_csv('web_scrapers/api-ninhas_dogs_data.csv')

# Define a function to find the closest matching breed name in the CSV file using fuzzy matching
def find_closest_match(dog_breed):
    closest_name = max(csv_df['name'], key=lambda name: fuzz.token_sort_ratio(dog_breed, name))
    return closest_name

# Apply the function to the 'dog_breed' column in the JSON dataframe to update the names in the file
json_df['dog_breed'] = json_df['dog_breed'].apply(find_closest_match)

# Save the updated JSON dataframe to a new file
json_df.to_json('web_scrapers/welton_reviews_cleaned.json', orient='records',  force_ascii=False, lines=True, indent=4)
