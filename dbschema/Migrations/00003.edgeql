CREATE MIGRATION m1kwv4eluyfzk6fj2yt5iy7wkh2x4tb7253cjwpfkcf3myojukuxoq
    ONTO m1leiewqjtptaz7i774zswvdkg7arbbjokuofcnmfhtawv4lsc4uia
{
  ALTER TYPE default::Breed {
      CREATE MULTI LINK groupID: default::Groups {
          CREATE PROPERTY id: std::uuid;
      };
      CREATE MULTI LINK imageID: default::Images {
          CREATE PROPERTY id: std::uuid;
      };
      CREATE LINK originID: default::Location {
          CREATE PROPERTY id: std::uuid;
      };
      CREATE MULTI LINK reviewID: default::Reviews {
          CREATE PROPERTY id: std::uuid;
      };
  };
};
