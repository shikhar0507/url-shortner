--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Ubuntu 13.3-1.pgdg20.04+1)
-- Dumped by pg_dump version 13.3 (Ubuntu 13.3-1.pgdg20.04+1)

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
-- Name: logs; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.logs (
    id integer NOT NULL,
    url character varying(200),
    username character varying(200),
    campaign character varying(200),
    source character varying(200),
    medium character varying(200),
    os character varying(20),
    browser character varying(20),
    device_type character varying(20),
    created_on timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.logs OWNER TO xanadu;

--
-- Data for Name: logs; Type: TABLE DATA; Schema: public; Owner: xanadu
--

COPY public.logs (id, url, username, campaign, source, medium, os, browser, device_type, created_on) FROM stdin;
1	https://hello.com	shikhar	summer_sale	facebook	digital	windows	chrome	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	facebook	digital	windows	chrome	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	facebook	digital	windows	chrome	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	facebook	digital	windows	firefox	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	facebook	digital	windows	firefox	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	facebook	digital	windows	safari	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	facebook	digital	linux	chrome	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	twitter	digital	linux	chrome	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	twitter	digital	linux	chrome	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	twitter	digital	linux	lynx	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	twitter	digital	linux	lynx	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	youtube	digital	osx	safari	desktop	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	youtube	digital	osx	safari	mobile	2021-05-21 00:20:27.232052+05:30
1	https://hello.com	shikhar	summer_sale	youtube	digital	linux	firefox	mobile	2021-05-21 00:20:27.232052+05:30
2	https://hello-new.com	shikhar	winter_sale	youtube	digital	linux	chrome	desktop	2021-05-21 00:20:27.232052+05:30
2	https://hello-new.com	shikhar	winter_sale	twitter	digital	linux	chrome	desktop	2021-05-21 00:20:27.232052+05:30
2	https://hello-new.com	shikhar	winter_sale	facebook	digital	linux	firefox	desktop	2021-05-21 00:20:27.232052+05:30
3	https://hello-just.com	shikhar_2	summer_sale	facebook	digital	linux	lynx	desktop	2021-05-21 00:30:06.484633+05:30
3	https://hello-just.com	shikhar_2	summer_sale	facebook	digital	linux	lynx	desktop	2021-05-21 00:30:06.484633+05:30
3	https://hello-just.com	shikhar_2	summer_sale	facebook	digital	linux	lynx	desktop	2021-05-21 00:30:06.484633+05:30
3	https://hello-just.com	shikhar_2	summer_sale	facebook	digital	linux	lynx	desktop	2021-05-21 00:30:06.484633+05:30
3	https://hello-just.com	shikhar_2	summer_sale	facebook	digital	linux	lynx	desktop	2021-05-21 00:30:06.484633+05:30
3	https://hello-just.com	shikhar_2	summer_sale	facebook	digital	linux	lynx	desktop	2021-05-21 00:30:06.484633+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
4	https://hello-just.com	shikhar_2	summer_sale	youtube	digital	linux	cli	desktop	2021-05-21 00:38:26.917382+05:30
\.


--
-- PostgreSQL database dump complete
--

