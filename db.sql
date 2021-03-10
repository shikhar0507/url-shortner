--
-- PostgreSQL database dump
--

-- Dumped from database version 13.0
-- Dumped by pg_dump version 13.0

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: urls; Type: TABLE; Schema: public; Owner: pujacapital
--

CREATE TABLE public.urls (
    id character varying(6) UNIQUE PRIMARY KEY,
    url character varying(2083)
);


ALTER TABLE public.urls OWNER TO pujacapital;

--
-- Data for Name: urls; Type: TABLE DATA; Schema: public; Owner: pujacapital
--

COPY public.urls (id, url) FROM stdin;
\.


--
-- PostgreSQL database dump complete
--
