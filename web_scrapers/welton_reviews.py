import requests
from bs4 import BeautifulSoup
import pandas as pd
import json

# Function to extract reviews from website:
def extract_dog_breed_reviews(base_url):
    response = requests.get(base_url)
    soup = BeautifulSoup(response.content, 'html.parser')

    dog_breed_reviews = []

    breed_links = [link['href'] for link in soup.select('.breeds a[href^="https://www.yourpurebredpuppy.com/reviews/"]')]

    for link in breed_links:
        breed_response = requests.get(link)
        breed_soup = BeautifulSoup(breed_response.text, 'html.parser')

        article_div = breed_soup.find('div', {'id': 'article'})

        if article_div:
            h1_tag = article_div.find('h1')
            dog_breed = h1_tag.text if h1_tag else None

            review_div = breed_soup.find('div', class_='article2')
            
            if review_div:
                paragraphs = review_div.find_all('p')
                review_text = ' '.join([p.get_text() for p in paragraphs])
            else:
                review_text = None

            dog_breed_reviews.append({'dog_breed': dog_breed, 'review': review_text})

    return dog_breed_reviews

# Function to create json file, removing "What's Good About 'Em, What's Bad About 'Em" text from page title
def create_json_file(reviews):
    for review in reviews:
        breed_name = review['dog_breed']
        if breed_name and ":" in breed_name:
            review['dog_breed'] = breed_name.split(":")[0].strip()
    return reviews

# Define base URL
base_url = "https://www.yourpurebredpuppy.com/dogbreeds/"
    
# Call function to extract dog breed reviews
dog_breed_reviews = extract_dog_breed_reviews(base_url)
    
# Call function to create DataFrame
welton_reviews = create_json_file(dog_breed_reviews)

# Drop duplicates from the json file
welton_reviews = [dict(t) for t in {tuple(d.items()) for d in welton_reviews}]

# Write data to a json file
with open('welton_reviews.json', 'w', encoding='utf-8') as json_file:
    json.dump(welton_reviews, json_file, ensure_ascii=False, indent=4)

#Display sample of JSON output
#sample_review = welton_reviews[:5]
#print(json.dumps(sample_review, ensure_ascii=False, indent=4))