--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Debian 16.6-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Homebrew)

-- Started on 2025-07-20 12:19:49 WITA

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
-- TOC entry 3374 (class 1262 OID 34859)
-- Name: userdb; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE userdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


\connect userdb

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
-- TOC entry 3375 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 34863)
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    fullname character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    code character varying DEFAULT ''::character varying NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);


--
-- TOC entry 216 (class 1259 OID 34878)
-- Name: users_auth; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users_auth (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying DEFAULT ''::character varying NOT NULL,
    phone_number character varying DEFAULT ''::character varying NOT NULL,
    password character varying DEFAULT ''::character varying NOT NULL,
    salt character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    user_uuid uuid NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);


--
-- TOC entry 3367 (class 0 OID 34863)
-- Dependencies: 215
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users VALUES ('baf64f79-0bc2-4333-b126-c36445e6851b', 'admin', '2025-07-12 02:19:13.556744+00', '2025-07-12 02:19:13.556744+00', 'U1', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');


--
-- TOC entry 3368 (class 0 OID 34878)
-- Dependencies: 216
-- Data for Name: users_auth; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users_auth VALUES ('1f7d4417-fdbc-46fd-a7a5-ddfecdfa0afa', 'admin@mail.com', '6285280042233', 'Pep/e5hZObYronYQeNj3c2BWAiU7RAZ0uG0m7SFEjDM=', 'rKxyNitzMfnfqmpWOsbezQ==', '2025-07-12 02:20:40.327997+00', '2025-07-12 02:20:40.327997+00', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');


--
-- TOC entry 3222 (class 2606 OID 34891)
-- Name: users_auth user_auth_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_auth
    ADD CONSTRAINT user_auth_pk PRIMARY KEY (uuid);


--
-- TOC entry 3220 (class 2606 OID 34874)
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk PRIMARY KEY (uuid);


--
-- TOC entry 3223 (class 1259 OID 34892)
-- Name: user_auth_user_uuid_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX user_auth_user_uuid_idx ON public.users_auth USING btree (user_uuid);


--
-- TOC entry 3218 (class 1259 OID 34875)
-- Name: users_code_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX users_code_idx ON public.users USING btree (code);


-- Completed on 2025-07-20 12:19:49 WITA

--
-- PostgreSQL database dump complete
--

