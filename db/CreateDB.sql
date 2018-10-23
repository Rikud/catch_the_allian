CREATE TABLE users (
  user_id BIGSERIAL PRIMARY KEY ,
  avatar character varying(255),
  email character varying(255),
  password character varying(255),
  user_name character varying(255)
);

CREATE TABLE history (
  id_history BIGSERIAL PRIMARY KEY,
  date_result TIMESTAMP,
  score integer,
  user_id bigint REFERENCES users(user_id)
);