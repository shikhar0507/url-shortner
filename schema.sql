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
-- Name: urls; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.urls (
    id character varying(6) NOT NULL,
    url character varying(2083),
    username character varying(20)
);


ALTER TABLE public.urls OWNER TO xanadu;

--
-- Name: log; Type: MATERIALIZED VIEW; Schema: public; Owner: xanadu
--

CREATE MATERIALIZED VIEW public.log AS
 WITH t AS (
         SELECT urls.id,
            urls.username,
            t2.device_type,
            t1.browser,
            t0.os,
            t4.total_clicks,
            sum(
                CASE
                    WHEN (t4.total_clicks IS NOT NULL) THEN t4.total_clicks
                    ELSE (0)::bigint
                END) AS sum
           FROM ((((public.urls
             LEFT JOIN ( SELECT logs.id,
                    count(logs.id) AS total_clicks
                   FROM public.logs
                  GROUP BY logs.id
                  ORDER BY (count(logs.id)) DESC) t4 ON ((t4.id = (urls.id)::text)))
             LEFT JOIN ( SELECT t_1.id,
                    t_1.device_type,
                    t_1.rank
                   FROM ( SELECT logs.id,
                            logs.device_type,
                            rank() OVER (PARTITION BY logs.id ORDER BY (count(logs.device_type)) DESC) AS rank
                           FROM public.logs
                          GROUP BY logs.id, logs.device_type) t_1
                  WHERE (t_1.rank = 1)) t2 ON (((urls.id)::text = t2.id)))
             LEFT JOIN ( SELECT t_1.id,
                    t_1.browser,
                    t_1.rank
                   FROM ( SELECT logs.id,
                            logs.browser,
                            rank() OVER (PARTITION BY logs.id ORDER BY (count(logs.browser)) DESC) AS rank
                           FROM public.logs
                          GROUP BY logs.id, logs.browser) t_1
                  WHERE (t_1.rank = 1)) t1 ON ((t2.id = t1.id)))
             LEFT JOIN ( SELECT t_1.id,
                    t_1.os,
                    t_1.rank
                   FROM ( SELECT logs.id,
                            logs.os,
                            rank() OVER (PARTITION BY logs.id ORDER BY (count(logs.os)) DESC) AS rank
                           FROM public.logs
                          GROUP BY logs.id, logs.os) t_1
                  WHERE (t_1.rank = 1)) t0 ON ((t1.id = t0.id)))
          GROUP BY urls.id, urls.username, t2.device_type, t1.browser, t0.os, t4.total_clicks
          ORDER BY (sum(
                CASE
                    WHEN (t4.total_clicks IS NOT NULL) THEN t4.total_clicks
                    ELSE (0)::bigint
                END)) DESC
        )
 SELECT t.id,
    t.username,
    t.browser,
    t.os,
    t.device_type,
    t.total_clicks
   FROM t
  WITH NO DATA;


ALTER TABLE public.log OWNER TO xanadu;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: xanadu
--

CREATE TABLE public.sessions (
    username character varying(20),
    sessionid character varying(100)
);


ALTER TABLE public.sessions OWNER TO xanadu;

--
-- Name: auth auth_pkey; Type: CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_pkey PRIMARY KEY (username);


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
-- Name: hash_idx_auth; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX hash_idx_auth ON public.auth USING hash (hash);


--
-- Name: id_idx_logs; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_logs ON public.logs USING btree (id);


--
-- Name: id_idx_urls; Type: INDEX; Schema: public; Owner: xanadu
--

CREATE INDEX id_idx_urls ON public.urls USING btree (id);


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
-- Name: logs url_fkey; Type: FK CONSTRAINT; Schema: public; Owner: xanadu
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT url_fkey FOREIGN KEY (id) REFERENCES public.urls(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

