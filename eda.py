import pandas as pd
from transformers import AutoTokenizer, AutoModelForTokenClassification
from transformers import pipeline

tokenizer = AutoTokenizer.from_pretrained("dslim/bert-base-NER")
model = AutoModelForTokenClassification.from_pretrained("dslim/bert-base-NER")

df_ninjas = pd.read_csv("web_scrapers/api-ninhas_dogs_data.csv")
df_ninjas.to_json("web_scrapers/api-ninhas_dogs_data.json")
df_dogtime = pd.read_json("web_scrapers/dogfacts.json")
df_ninjas.loc[df_ninjas.name == "Irish Red and White Setter", "name"] = "Irish Red And White Setter"
df_ninjas.index = df_ninjas.name
del df_ninjas["name"]
df_dogfact = df_ninjas.join(df_dogtime, how="outer").drop_duplicates()
df_dogfact.to_json('web_scrapers/dogfacts.json')
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
