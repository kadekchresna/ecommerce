--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Debian 16.6-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Homebrew)

-- Started on 2025-07-15 00:36:26 WITA

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
-- TOC entry 3363 (class 1262 OID 35417)
-- Name: shopdb; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE shopdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


\connect shopdb

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
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


--
-- TOC entry 3364 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 843 (class 1247 OID 35433)
-- Name: outbox_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.outbox_status_type AS ENUM (
    'created',
    'in-progress',
    'failed',
    'success'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 35418)
-- Name: shops; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shops (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    code character varying DEFAULT ''::character varying NOT NULL,
    user_uuid uuid NOT NULL,
    name character varying DEFAULT ''::character varying NOT NULL,
    "desc" text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);


--
-- TOC entry 3357 (class 0 OID 35418)
-- Dependencies: 215
-- Data for Name: shops; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.shops (uuid, code, user_uuid, name, "desc", created_at, updated_at, created_by, updated_by) FROM stdin;
ee7ae731-2eed-42f6-8852-d84c8847f112	gramedia-1	baf64f79-0bc2-4333-b126-c36445e6851b	Gramedia	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas finibus diam et arcu ultrices laoreet. Nulla et quam fermentum, tincidunt mauris nec, vestibulum diam. Nulla pulvinar diam ipsum, gravida porttitor nulla finibus sit amet. Sed faucibus, sapien vitae ornare tempus, nunc nisl rhoncus magna, quis imperdiet ligula lorem in nulla. Fusce dapibus quis lorem consectetur viverra. Phasellus maximus purus nec ligula iaculis, ac aliquet tortor fringilla. Cras imperdiet id magna ac convallis.\n\n	2025-07-13 18:45:11.235314+00	2025-07-13 18:45:11.235314+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
\.


--
-- TOC entry 3212 (class 2606 OID 35430)
-- Name: shops shop_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shops
    ADD CONSTRAINT shop_pk PRIMARY KEY (uuid);


--
-- TOC entry 3213 (class 1259 OID 35431)
-- Name: shops_code_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX shops_code_idx ON public.shops USING btree (code);


-- Completed on 2025-07-15 00:36:26 WITA

--
-- PostgreSQL database dump complete
--

