import requests
import re
import pandas as pd
import numpy as np
from bs4 import BeautifulSoup


class DogBreed(object):
    def __init__(self, url):
        self.url = url
        self.my_breeds = {}
        self.df_facts = pd.DataFrame()

    def stripHTML(self, data):
        p = re.compile(r'[\t\n]')
        return p.sub('', str(data))

    def breeds(self):
        # downloaded the dynamic page to get all the dog breeds
        with open('dogtime_breeds.html', 'r') as htmlFile:
            soup = BeautifulSoup(htmlFile.read(), 'html.parser')
        for i, node in enumerate(soup.find_all('div', {'class': 'wp-block-xwp-curated-content__card-header'})):
            a = node.find_all(href=True)
            k = self.stripHTML(a[0].text)
            v = a[0]['href']
            if self.url in v:
                self.my_breeds[k] = v

    def facts(self):
        for breedName in self.my_breeds.keys():
            breedurl = self.my_breeds[breedName]
            response = requests.get(breedurl)
            soup = BeautifulSoup(response.text, 'html.parser')
            mydict = {}
            for i, node in enumerate(soup.find_all('li', {'class': ''})):
                # we only need 12 facts
                if i == 12:
                    break
                try:
                    k, v = node.text.split(':', 1)
                    v = v.strip()
                    mydict[k] = v
                except:
                    pass
            self.df_facts = pd.concat([self.df_facts, pd.DataFrame(mydict, index=[breedName])])


if __name__ == "__main__":
    url = "https://dogtime.com/dog-breeds/"
    dog = DogBreed(url)
    dog.breeds()
    dog.facts()
    dog.df_facts.dropna(how='all').to_json('dogfacts.json')
