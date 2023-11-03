CREATE MIGRATION m1e35o7n4jtjhygfqmq5mbteb63fhe6ybacjg3k4mien7t4u3xj6ua
    ONTO initial
{
  CREATE EXTENSION pgvector VERSION '0.4';
  CREATE TYPE default::Breed {
      CREATE PROPERTY barking: std::int32;
      CREATE PROPERTY coat_length: std::int32;
      CREATE PROPERTY drooling: std::int32;
      CREATE PROPERTY energy: std::int32;
      CREATE PROPERTY good_with_children: std::int32;
      CREATE PROPERTY good_with_other_dogs: std::int32;
      CREATE PROPERTY good_with_strangers: std::int32;
      CREATE PROPERTY grooming: std::int32;
      CREATE PROPERTY height: std::decimal;
      CREATE PROPERTY life_expectancy: std::decimal;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY playfulness: std::int32;
      CREATE PROPERTY protectiveness: std::int32;
      CREATE PROPERTY shedding: std::int32;
      CREATE PROPERTY trainability: std::int32;
      CREATE PROPERTY weight: std::decimal;
  };
  CREATE TYPE default::Images {
      CREATE LINK BreedFrom: default::Breed {
          CREATE PROPERTY name: std::str;
      };
      CREATE PROPERTY image_link: std::str;
  };
  CREATE TYPE default::Location {
      CREATE LINK BreedFrom: default::Breed {
          CREATE PROPERTY name: std::str;
      };
      CREATE PROPERTY City: std::str;
      CREATE PROPERTY Country: std::str;
      CREATE PROPERTY State: std::str;
  };
  CREATE SCALAR TYPE default::TxEmbedding EXTENDING ext::pgvector::vector<768>;
  CREATE TYPE default::Reviews {
      CREATE LINK BreedFrom: default::Breed {
          CREATE PROPERTY name: std::str;
      };
      CREATE REQUIRED PROPERTY embedding: default::TxEmbedding;
      CREATE PROPERTY review: std::str;
  };
  CREATE TYPE default::Groups {
      CREATE MULTI LINK breeds_list: default::Breed {
          CREATE PROPERTY name: std::str;
      };
      CREATE PROPERTY group_description: std::str;
      CREATE REQUIRED PROPERTY group_name: std::str;
  };
};
