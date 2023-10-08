import requests
from bs4 import BeautifulSoup
import pandas as pd

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

# Function to create dataframe, removing "What's Good About 'Em, What's Bad About 'Em" text from page title
def create_dataframe(reviews):
    for review in reviews:
        breed_name = review['dog_breed']
        if breed_name and ":" in breed_name:
            review['dog_breed'] = breed_name.split(":")[0].strip()
    return pd.DataFrame(reviews)

# Define base URL
base_url = "https://www.yourpurebredpuppy.com/dogbreeds/"
    
# Call function to extract dog breed reviews
dog_breed_reviews = extract_dog_breed_reviews(base_url)
    
# Call function to create DataFrame
welton_reviews_df = create_dataframe(dog_breed_reviews)

# Drop duplicates from the dataframe
welton_reviews_df.drop_duplicates(inplace=True)

#Optional - display dataframe
#welton_reviews_df

#Optional - write dataframe to csv
#welton_reviews_df.to_csv('welton_reviews.csv', index=False)
