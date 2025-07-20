--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Debian 16.6-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Homebrew)

-- Started on 2025-07-15 00:35:14 WITA

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
-- TOC entry 3424 (class 1262 OID 34861)
-- Name: omsdb; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE omsdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


\connect omsdb

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
-- TOC entry 3425 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 861 (class 1247 OID 35476)
-- Name: inbox_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.inbox_status_type AS ENUM (
    'created',
    'in-progress',
    'failed',
    'success'
);


--
-- TOC entry 843 (class 1247 OID 34910)
-- Name: order_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.order_status_type AS ENUM (
    'created',
    'in-payment',
    'expired',
    'completed',
    'cancelled',
    'reserving-stock',
    'failed'
);


--
-- TOC entry 858 (class 1247 OID 35455)
-- Name: outbox_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.outbox_status_type AS ENUM (
    'created',
    'in-progress',
    'failed',
    'success'
);


--
-- TOC entry 852 (class 1247 OID 34983)
-- Name: warehouse_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.warehouse_status_type AS ENUM (
    'active',
    'inactive'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 35485)
-- Name: inbox; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inbox (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    status public.inbox_status_type DEFAULT 'created'::public.inbox_status_type NOT NULL,
    type character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    reference character varying DEFAULT ''::character varying NOT NULL,
    response jsonb DEFAULT '{}'::jsonb NOT NULL,
    retry_count integer DEFAULT 0 NOT NULL,
    action character varying DEFAULT ''::character varying NOT NULL
);


--
-- TOC entry 215 (class 1259 OID 34919)
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    code character varying DEFAULT ''::character varying NOT NULL,
    user_uuid uuid NOT NULL,
    total_amount double precision DEFAULT 0.0 NOT NULL,
    expired_at timestamp with time zone NOT NULL,
    status public.order_status_type DEFAULT 'created'::public.order_status_type NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL
);


--
-- TOC entry 216 (class 1259 OID 34932)
-- Name: orders_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders_detail (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    product_uuid uuid NOT NULL,
    product_title character varying DEFAULT ''::character varying NOT NULL,
    product_price double precision DEFAULT 0.0 NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    sub_total double precision DEFAULT 0.0 NOT NULL,
    order_uuid uuid NOT NULL
);


--
-- TOC entry 217 (class 1259 OID 35441)
-- Name: outbox; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.outbox (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    status public.outbox_status_type DEFAULT 'created'::public.outbox_status_type NOT NULL,
    type character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    reference character varying DEFAULT ''::character varying NOT NULL,
    response jsonb DEFAULT '{}'::jsonb NOT NULL,
    retry_count integer DEFAULT 0 NOT NULL,
    action character varying DEFAULT ''::character varying NOT NULL
);


--
-- TOC entry 3268 (class 2606 OID 35717)
-- Name: inbox inbox_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inbox
    ADD CONSTRAINT inbox_pk PRIMARY KEY (uuid);


--
-- TOC entry 3271 (class 2606 OID 35719)
-- Name: inbox inbox_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inbox
    ADD CONSTRAINT inbox_unique UNIQUE (type, reference, action);


--
-- TOC entry 3260 (class 2606 OID 34931)
-- Name: orders orders_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pk PRIMARY KEY (uuid);


--
-- TOC entry 3263 (class 2606 OID 35452)
-- Name: outbox outbox_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.outbox
    ADD CONSTRAINT outbox_pk PRIMARY KEY (uuid);


--
-- TOC entry 3266 (class 2606 OID 35721)
-- Name: outbox outbox_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.outbox
    ADD CONSTRAINT outbox_unique UNIQUE (type, reference, action);


--
-- TOC entry 3269 (class 1259 OID 35534)
-- Name: inbox_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX inbox_status_idx ON public.inbox USING btree (status);


--
-- TOC entry 3258 (class 1259 OID 35011)
-- Name: orders_code_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX orders_code_idx ON public.orders USING btree (code);


--
-- TOC entry 3261 (class 1259 OID 34942)
-- Name: orders_detail_order_uuid_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX orders_detail_order_uuid_idx ON public.orders_detail USING btree (order_uuid);


--
-- TOC entry 3264 (class 1259 OID 35525)
-- Name: outbox_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX outbox_status_idx ON public.outbox USING btree (status);


-- Completed on 2025-07-15 00:35:14 WITA

--
-- PostgreSQL database dump complete
--

