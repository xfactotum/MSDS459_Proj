CREATE MIGRATION m1leiewqjtptaz7i774zswvdkg7arbbjokuofcnmfhtawv4lsc4uia
    ONTO m1e35o7n4jtjhygfqmq5mbteb63fhe6ybacjg3k4mien7t4u3xj6ua
{
  ALTER TYPE default::Images {
      ALTER LINK BreedFrom {
          RENAME TO BreedImage;
      };
  };
  ALTER TYPE default::Location {
      ALTER PROPERTY City {
          RENAME TO city;
      };
  };
  ALTER TYPE default::Location {
      ALTER PROPERTY Country {
          RENAME TO country;
      };
  };
  ALTER TYPE default::Location {
      ALTER PROPERTY State {
          RENAME TO state;
      };
  };
  ALTER TYPE default::Reviews {
      ALTER LINK BreedFrom {
          RENAME TO BreedReview;
      };
  };
};
