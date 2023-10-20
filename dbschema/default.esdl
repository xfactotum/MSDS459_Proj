using extension pgvector;

module default {
    scalar type TxEmbedding extending
    ext::pgvector::vector<1536>;

    type Breed {
        required property name -> str {
            constraint exclusive
        };
        property height -> str; #akc extraction
        property weight -> str; #akc extraction
        property life_expectancy -> str; #akc extraction
        property breed_group -> str; #akc extraction
        property good_with_children -> int32; #api extraction
        property good_with_other_dogs -> int32; #api extraction
        property shedding -> int32; #api extraction
        property grooming -> int32; #api extraction
        property drooling -> int32; #api extraction
        property coat_length -> int32; #api extraction
        property good_with_strangers -> int32; #api extraction
        property playfulness -> int32; #api extraction
        property protectiveness -> int32; #api extraction
        property trainability -> int32; #api extraction
        property energy -> int32; #api extraction
        property barking -> int32; #api extraction
        property origin -> int32; #dogtime extraction

        multi link GroupList -> Groups;
        
    };

    type Groups {
        required property group_name -> str;
        property description -> str; #akc website - manually create records?
        link BreedGroup -> Breed {
            property breed_group ->str;
        };
    }

    type Reviews {
        required property name -> str;
        property review -> str;
        link BreedReview -> Breed;
        required embedding: TxEmbedding;
    }

    type Location {
        required property name -> str;
        property City -> str;
        property State -> str;
        property Country -> str;
        link BreedFrom -> Breed;
    }

    type Images {
        required property name -> str;
        property image_link -> str; #api extraction
        link BreedPicture -> Breed;
    }

}
