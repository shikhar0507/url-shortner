--
-- PostgreSQL database dump
--

-- Dumped from database version 11.12 (Debian 11.12-0+deb10u1)
-- Dumped by pg_dump version 11.12 (Debian 11.12-0+deb10u1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: insertlongurl(text); Type: FUNCTION; Schema: public; Owner: xanadu
--

CREATE FUNCTION public.insertlongurl(url text) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
       DECLARE
	nextId INTEGER;
	indexMapping TEXT[] := '{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","0","1","2","3","4","5","6","7","8","9"}';
	modval INTEGER;
	S TEXT;
	
       BEGIN
        loop
	--	BEGIN		
				SELECT  last_value + CASE WHEN is_called THEN 1 ELSE 0 END FROM urls_seq_seq INTO nextId;
				while nextId > 0 loop
	      	      		      modval := nextId % 62;
	      	 		      nextId := nextId / 62;
	      	    	 	      SELECT CONCAT(indexMapping[modval],S) INTO S;
			        end loop;
				BEGIN
					RAISE NOTICE 'mil';
					INSERT INTO urls(id) VALUES(S);
					RETURN S;
				EXCEPTION WHEN unique_violation THEN
			  		  -- do nothing
					  	  
					  RAISE notice  'duplicate id %', S;
				END;
	--	COMMIT;
	END LOOP;
       END;
$$;


ALTER FUNCTION public.insertlongurl(url text) OWNER TO xanadu;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: auth; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.auth (
    username character varying(20) NOT NULL,
    hash character varying(200) NOT NULL
);


ALTER TABLE public.auth OWNER TO xanadu;

--
-- Name: country; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.country (
    id integer NOT NULL,
    name character varying(80) NOT NULL
);


ALTER TABLE public.country OWNER TO xanadu;

--
-- Name: country_block; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.country_block (
    id character varying(6) NOT NULL,
    name character varying(80) NOT NULL
);


ALTER TABLE public.country_block OWNER TO xanadu;

--
-- Name: expiration; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.expiration (
    id character varying(6) NOT NULL,
    expiration timestamp with time zone NOT NULL,
    expired_url text DEFAULT 'http://localhost:8080'::text
);


ALTER TABLE public.expiration OWNER TO xanadu;

--
-- Name: ip_address_to_country; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.ip_address_to_country (
    ip cidr NOT NULL,
    country character varying(80) NOT NULL
);


ALTER TABLE public.ip_address_to_country OWNER TO xanadu;

--
-- Name: ip_redirect; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.ip_redirect (
    id character varying(6) NOT NULL,
    url text NOT NULL,
    country character varying(80) NOT NULL,
    redirect_url text NOT NULL
);


ALTER TABLE public.ip_redirect OWNER TO xanadu;

--
-- Name: logs; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.logs (
    id text NOT NULL,
    campaign character varying(200),
    source character varying(200),
    medium character varying(200),
    os character varying(20),
    browser character varying(20),
    device_type character varying(20),
    created_on timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ip cidr
)
WITH (parallel_workers='6');


ALTER TABLE public.logs OWNER TO xanadu;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.sessions (
    username character varying(20),
    sessionid character varying(100)
);


ALTER TABLE public.sessions OWNER TO xanadu;

--
-- Name: urls; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.urls (
    id character varying(6) NOT NULL,
    url character varying(2083),
    username character varying(20),
    tag character varying(20),
    password character varying(200),
    not_found_url text,
    android_deep_link text,
    ios_deep_link text,
    seq integer NOT NULL
);


ALTER TABLE public.urls OWNER TO xanadu;

--
-- Name: urls_seq_seq; Type: SEQUENCE; Schema: public; Owner: xanadu
--

CREATE SEQUENCE public.urls_seq_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.urls_seq_seq OWNER TO xanadu;

--
-- Name: urls_seq_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xanadu
--

ALTER SEQUENCE public.urls_seq_seq OWNED BY public.urls.seq;


--
-- Name: urls seq; Type: DEFAULT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.urls ALTER COLUMN seq SET DEFAULT nextval('public.urls_seq_seq'::regclass);


--
-- Name: auth auth_pkey; Type: CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_pkey PRIMARY KEY (username);


--
-- Name: country country_pkey; Type: CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_pkey PRIMARY KEY (id);


--
-- Name: expiration expiration_pkey; Type: CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.expiration
    ADD CONSTRAINT expiration_pkey PRIMARY KEY (id);


--
-- Name: auth u_hash; Type: CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT u_hash UNIQUE (hash);


--
-- Name: urls urls_pkey; Type: CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.urls
    ADD CONSTRAINT urls_pkey PRIMARY KEY (id);


--
-- Name: country_idx_ip_redirect; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX country_idx_ip_redirect ON public.ip_redirect USING btree (country);


--
-- Name: hash_idx_auth; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX hash_idx_auth ON public.auth USING hash (hash);


--
-- Name: id_idx_country_block; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_country_block ON public.country_block USING btree (id);


--
-- Name: id_idx_expiration; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_expiration ON public.expiration USING btree (id);


--
-- Name: id_idx_ip_redirect; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_ip_redirect ON public.ip_redirect USING btree (id);


--
-- Name: id_idx_logs; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_logs ON public.logs USING btree (id);


--
-- Name: id_idx_urls; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_urls ON public.urls USING btree (id);


--
-- Name: ip_idx_ip_address_to_country; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX ip_idx_ip_address_to_country ON public.ip_address_to_country USING btree (ip);


--
-- Name: name_idx_country; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX name_idx_country ON public.country USING btree (name);


--
-- Name: name_idx_country_block; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX name_idx_country_block ON public.country_block USING btree (name);


--
-- Name: sessionid_idx_sessions; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX sessionid_idx_sessions ON public.sessions USING hash (username);


--
-- Name: username_idx_auth; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX username_idx_auth ON public.auth USING btree (username);


--
-- Name: username_idx_sessions; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX username_idx_sessions ON public.sessions USING btree (username);


--
-- Name: username_idx_urls; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX username_idx_urls ON public.urls USING btree (username);


--
-- Name: expiration foreign_fk_id; Type: FK CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.expiration
    ADD CONSTRAINT foreign_fk_id FOREIGN KEY (id) REFERENCES public.urls(id) ON DELETE CASCADE;


--
-- Name: logs url_fkey; Type: FK CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT url_fkey FOREIGN KEY (id) REFERENCES public.urls(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

