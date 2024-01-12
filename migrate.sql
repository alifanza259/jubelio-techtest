-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	email varchar NOT NULL,
	"password" varchar NOT NULL,
	"name" varchar NOT NULL,
	image_url varchar NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- public.messages definition

-- Drop table

-- DROP TABLE public.messages;

CREATE TABLE public.messages (
	id serial4 NOT NULL,
	sender_id int4 NOT NULL,
	receiver_id int4 NOT NULL,
	"content" text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT messages_pkey PRIMARY KEY (id)
);

INSERT INTO public.users (email,"password","name",image_url,created_at) VALUES
	 ('a@gmail.com','$2a$12$S0sss.CJrY3IQh2erifDveDQK5TdM2/9YfMDuwHOuPdoH.8OH7KyO','Andi','https://bucket-alif.s3.ap-southeast-1.amazonaws.com/angka1.jpeg','2024-01-11 18:10:03.153'),
	 ('b@gmail.com','$2a$12$S0sss.CJrY3IQh2erifDveDQK5TdM2/9YfMDuwHOuPdoH.8OH7KyO','Bakri','https://bucket-alif.s3.ap-southeast-1.amazonaws.com/angka2.png','2024-01-11 18:10:03.153'),
	 ('c@gmail.com','$2a$12$S0sss.CJrY3IQh2erifDveDQK5TdM2/9YfMDuwHOuPdoH.8OH7KyO','Ciko','https://bucket-alif.s3.ap-southeast-1.amazonaws.com/angka3.png','2024-01-11 18:10:03.153');
