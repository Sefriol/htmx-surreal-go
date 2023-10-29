DEFINE TABLE relative SCHEMAFULL;
DEFINE FIELD Parent ON TABLE relatives TYPE record<user>;
DEFINE FIELD Child ON TABLE relatives TYPE record<user>;
DEFINE INDEX unique_relationships 
    ON TABLE relative 
    COLUMNS in, out UNIQUE;
