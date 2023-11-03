import json

import pandas as pd
import edgedb

client = edgedb.create_client()

def breedLoader():
    # df_breed = pd.read_json("web_scrapers/breed.json")
    # df_breed.to_json("data_files/breedToLoad.json", orient='records')
    with open("data_files/breedToLoad.json", "r") as f:
        jsonBreed = json.load(f)
        jsonBreeds = json.dumps(jsonBreed)

    client.query("""
        with
            raw_data :=  <json>$data,
        for item in json_array_unpack(raw_data) union (
            insert Breed { 
                name := <str>item['name'],
                height := <decimal>item['height'],
                weight := <decimal>item['weight'],
                life_expectancy := <decimal>item['life_expectancy'],
                good_with_children := <int32>item['good_with_children'],
                good_with_other_dogs := <int32>item['good_with_other_dogs'],
                good_with_strangers := <int32>item['good_with_strangers'],
                shedding := <int32>item['shedding'],
                grooming := <int32>item['grooming'],
                drooling := <int32>item['drooling'],
                coat_length := <int32>item['coat_length'],
                playfulness := <int32>item['playfulness'],
                protectiveness := <int32>item['protectiveness'],
                trainability := <int32>item['trainability'],
                energy := <int32>item['energy'],
                barking := <int32>item['barking']            
            }
        );
    """, data=jsonBreeds)

    result = client.query("""
        select Breed {name}
    """)


def reviewLoader():
    # df_review = pd.read_json("web_scrapers/reviews_with_embeddings.json")
    # df_review.to_json("data_files/reviewsToLoad.json", orient='records')
    with open("data_files/reviewsToLoad.json", "r") as f:
        jsonReview = json.load(f)
        jsonReviews = json.dumps(jsonReview)

    client.query("""
        with
            raw_data :=  <json>$data,
        for item in json_array_unpack(raw_data) union (
            insert Reviews { 
                BreedFrom := <Breed>(
                    select detached Breed
                    filter .name = <str>item['name']
                    ),
                review := <str>item['review'],
                embedding := <TxEmbedding>item['embedding']
            }
        );
    """, data=jsonReviews)


def locLoader():
    # df_location = pd.read_json("web_scrapers/location.json")
    # del df_location["Unnamed: 0"]
    # del df_location["id"]
    # df_location.columns = ['name', 'country', 'city', 'state']
    # df_location.to_json("data_files/locationToLoad.json", orient='records')
    with open("data_files/locationToLoad.json", "r") as f:
        jsonLocation = json.load(f)
        jsonLocations = json.dumps(jsonLocation)

    client.query("""
        with
            raw_data :=  <json>$data,
        for item in json_array_unpack(raw_data) union (
            insert Location { 
                BreedFrom := <Breed>(
                    select detached Breed
                    filter .name = <str>item['name']
                    ),
                Country := <str>item['country'],
                City := <str>item['city'],
                State := <str>item['state']
            }
        );
    """, data=jsonLocations)


def imageLoader():
    # df_image = pd.read_json("web_scrapers/images.json")
    # df_image.to_json("data_files/imageToLoad.json", orient='records')
    with open("data_files/imageToLoad.json", "r") as f:
        jsonImage = json.load(f)
        jsonImages = json.dumps(jsonImage)

    client.query("""
        with
            raw_data :=  <json>$data,
        for item in json_array_unpack(raw_data) union (
            insert Images { 
                BreedFrom := <Breed>(
                    select detached Breed
                    filter .name = <str>item['name']
                    ),
                image_link := <str>item['image_link']
            }
        );
    """, data=jsonImages)


def groupLoader():
    # df_group = pd.read_json("web_scrapers/akc_groups.json")
    # df_group.to_json("data_files/groupToLoad.json", orient='records')
    with open("data_files/groupToLoad.json", "r") as f:
        jsonGroup = json.load(f)
        jsonGroups = json.dumps(jsonGroup)

    client.query("""
        with
            raw_data :=  <json>$data,
        for item in json_array_unpack(raw_data) union (
            with breeds := (
                select Breed filter .name in array_unpack(<array<str>>item['breeds_list'])
            )
            insert Groups { 
                breeds_list := breeds,
                group_name := <str>item['group_name'],
                group_description := <str>item['group_description']
            }
        );
    """, data=jsonGroups)

breedLoader()
reviewLoader()
locLoader()
imageLoader()
groupLoader()
