CREATE MIGRATION m1jstmmmpl7gu6645xzx57uizxxd7atgcfrzledauqkb3g3fif6raq
    ONTO initial
{
  CREATE EXTENSION pgvector VERSION '0.4';
  CREATE TYPE default::Breed {
      CREATE PROPERTY barking: std::int32;
      CREATE PROPERTY breed_group: std::str;
      CREATE PROPERTY coat_length: std::int32;
      CREATE PROPERTY drooling: std::int32;
      CREATE PROPERTY energy: std::int32;
      CREATE PROPERTY good_with_children: std::int32;
      CREATE PROPERTY good_with_other_dogs: std::int32;
      CREATE PROPERTY good_with_strangers: std::int32;
      CREATE PROPERTY grooming: std::int32;
      CREATE PROPERTY height: std::str;
      CREATE PROPERTY life_expectancy: std::str;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY origin: std::int32;
      CREATE PROPERTY playfulness: std::int32;
      CREATE PROPERTY protectiveness: std::int32;
      CREATE PROPERTY shedding: std::int32;
      CREATE PROPERTY trainability: std::int32;
      CREATE PROPERTY weight: std::str;
  };
  CREATE TYPE default::Groups {
      CREATE LINK BreedGroup: default::Breed {
          CREATE PROPERTY breed_group: std::str;
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY group_name: std::str;
  };
  ALTER TYPE default::Breed {
      CREATE MULTI LINK GroupList: default::Groups;
  };
  CREATE TYPE default::Location {
      CREATE LINK BreedFrom: default::Breed;
      CREATE PROPERTY City: std::str;
      CREATE PROPERTY Country: std::str;
      CREATE PROPERTY State: std::str;
      CREATE REQUIRED PROPERTY name: std::str;
  };
  CREATE TYPE default::Images {
      CREATE LINK BreedPicture: default::Breed;
      CREATE PROPERTY image_link: std::str;
      CREATE REQUIRED PROPERTY name: std::str;
  };
  CREATE SCALAR TYPE default::TxEmbedding EXTENDING ext::pgvector::vector<1536>;
  CREATE TYPE default::Reviews {
      CREATE LINK BreedReview: default::Breed;
      CREATE REQUIRED PROPERTY embedding: default::TxEmbedding;
      CREATE REQUIRED PROPERTY name: std::str;
      CREATE PROPERTY review: std::str;
  };
};
