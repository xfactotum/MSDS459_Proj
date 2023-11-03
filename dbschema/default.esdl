using extension pgvector;

module default {
    scalar type TxEmbedding extending
    ext::pgvector::vector<768>;

    type Breed {
        required property name -> str { #api extraction
            constraint exclusive
        };
        property height -> decimal; #api extraction
        property weight -> decimal; #api extraction
        property life_expectancy -> decimal; #api extraction
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
    };

    type Reviews {
        link BreedFrom -> Breed {
            property name -> str;
        };
        property review -> str; #reviews extraction
        required embedding: TxEmbedding;
    }

    type Location {
        link BreedFrom -> Breed {
            property name -> str;
        };
        property City -> str; #dogtime extraction
        property State -> str; #dogtime extraction
        property Country -> str; #dogtime extraction
    }

    type Images {
        link BreedFrom -> Breed {
            property name -> str;
        };
        property image_link -> str; #api extraction
    }

    type Groups {
        required property group_name -> str; #akc extraction
        property group_description -> str; #akc extraction
        multi link breeds_list -> Breed {
            property name ->str;
        };
    }

}
