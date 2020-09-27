USE balance;

CREATE TABLE user
(
  id              INT(20) unsigned NOT NULL AUTO_INCREMENT,
  name            VARCHAR(255) NOT NULL,       
  PRIMARY KEY     (id)                               
);

CREATE TABLE balance
(
  id              INT(20) unsigned NOT NULL AUTO_INCREMENT,
  uid             INT(20) unsigned NOT NULL,
  value           FLOAT(20) NOT NULL,     
  PRIMARY KEY     (id)                               
);

CREATE TABLE note
(
  id              INT(20) unsigned NOT NULL AUTO_INCREMENT,
  bid             INT(20) unsigned NOT NULL,
  text            VARCHAR(255) NOT NULL,   
  timestamp       INT(20) NOT NULL,     
  PRIMARY KEY     (id)                               
);

