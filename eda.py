import pandas as pd
from sentence_transformers import SentenceTransformer
from transformers import pipeline
from sentence_transformers import SentenceTransformer

# Load model from HuggingFace Hub
tokenizer = AutoTokenizer.from_pretrained("dslim/bert-base-NER")
model = AutoModelForTokenClassification.from_pretrained("dslim/bert-base-NER")


def breed():
    # df_ninjas = pd.read_csv("web_scrapers/api-ninhas_dogs_data.csv")
    # df_ninjas.to_json("web_scrapers/api-ninhas_dogs_data.json")
    # df_dogtime = pd.read_json("web_scrapers/dogfacts.json")
    # df_ninjas.loc[df_ninjas.name == "Irish Red and White Setter", "name"] = "Irish Red And White Setter"
    # df_ninjas.index = df_ninjas.name
    # del df_ninjas["name"]
    df_avg = pd.read_csv("web_scrapers/dog_avg_stat.csv")
    df_avg.columns = ["name","life_expectancy","height","weight"]
    df_dogtime = pd.read_csv("web_scrapers/dogfacts.csv")
    df_dogtime = df_dogtime[[
        "dog_breed",
        "good_with_children",
        "good_with_other_dogs",
        "shedding",
        "grooming",
        "drooling",
        "coat_length",
        "good_with_strangers",
        "playfulness",
        "protectiveness",
        "trainability",
        "energy",
        "barking"]]
    df_dogtime.columns = [
        "name",
        "good_with_children",
        "good_with_other_dogs",
        "shedding",
        "grooming",
        "drooling",
        "coat_length",
        "good_with_strangers",
        "playfulness",
        "protectiveness",
        "trainability",
        "energy",
        "barking"]
    df_dogfact = df_avg.merge(df_dogtime, how="outer", on="name").drop_duplicates()
    df_dogfact.to_csv('web_scrapers/breed.csv')
    df_dogfact.to_json('web_scrapers/breed.json')

def getEmbeddings():
    model = SentenceTransformer('sentence-transformers/all-mpnet-base-v2')
    df_review = pd.read_json("web_scrapers/reviews_cleaned.json")
    df_review["embedding"] = ""
    for i, row in df_review.iterrows():
        sentences = row["review"]
        embedding = model.encode(sentences)
        df_review.loc[i, "embedding"] = embedding

    df_review.to_json("web_scrapers/reviews_with_embeddings.json")


def getLocation():
    df_loc = pd.DataFrame()
    nlp = pipeline("ner", model=model, tokenizer=tokenizer)
    for pk, breed in df_dogfact[~df_dogfact['Origin'].isnull()][['Origin']].iterrows():
        splitOrig = breed.Origin.split(',')
        ner_results = nlp(breed.Origin)
        locNER = [x['word'] for x in ner_results if '-LOC' in x['entity']]
        if (len(locNER) > len(splitOrig)) & (len(splitOrig) > 2):
            df_loc = pd.concat([df_loc,
                                pd.DataFrame({
                                    "id": [df_loc.shape[0]],
                                    "Breed": [pk],
                                    "City": [' '.join(splitOrig[:-2])],
                                    "State": [splitOrig[-2]],
                                    "Country": [splitOrig[-1]]})],
                               ignore_index=True)
        else:
            myloc = " ".join(locNER)
            if myloc != '':
                df_loc = pd.concat([df_loc,
                                    pd.DataFrame({"id": [df_loc.shape[0]], "Breed": [pk], "Country": [myloc]})],
                                   ignore_index=True)
    df_loc.to_json('web_scrapers/location.json')


getEmbeddings()
# breed()
# getLocation()
