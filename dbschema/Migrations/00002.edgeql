CREATE MIGRATION m1njwm6td4riatjntp6lmunihd7lfqr5skxq32dsk4zd774zsxfdgq
    ONTO m1jstmmmpl7gu6645xzx57uizxxd7atgcfrzledauqkb3g3fif6raq
{
  ALTER TYPE default::Breed {
      DROP PROPERTY breed_group;
  };
  ALTER TYPE default::Groups {
      CREATE PROPERTY breeds_list: array<std::str>;
  };
  ALTER TYPE default::Groups {
      ALTER PROPERTY description {
          RENAME TO group_description;
      };
  };
  ALTER TYPE default::Location {
      ALTER LINK BreedFrom {
          CREATE PROPERTY origin: std::str;
      };
  };
};
