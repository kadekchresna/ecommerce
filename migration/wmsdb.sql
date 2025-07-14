--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Debian 16.6-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Homebrew)

-- Started on 2025-07-15 00:37:26 WITA

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
-- TOC entry 3423 (class 1262 OID 34862)
-- Name: wmsdb; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE wmsdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


\connect wmsdb

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
-- TOC entry 3424 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 855 (class 1247 OID 35604)
-- Name: inbox_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.inbox_status_type AS ENUM (
    'created',
    'in-progress',
    'failed',
    'success'
);


--
-- TOC entry 852 (class 1247 OID 35595)
-- Name: outbox_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.outbox_status_type AS ENUM (
    'created',
    'in-progress',
    'failed',
    'success'
);


--
-- TOC entry 849 (class 1247 OID 34996)
-- Name: warehouse_status_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.warehouse_status_type AS ENUM (
    'active',
    'inactive'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 35613)
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
-- TOC entry 218 (class 1259 OID 35629)
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
-- TOC entry 215 (class 1259 OID 34957)
-- Name: warehouses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouses (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying DEFAULT ''::character varying NOT NULL,
    code character varying DEFAULT ''::character varying NOT NULL,
    "desc" text NOT NULL,
    shop_uuid uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    status public.warehouse_status_type DEFAULT 'active'::public.warehouse_status_type NOT NULL
);


--
-- TOC entry 216 (class 1259 OID 34971)
-- Name: warehouses_stock; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouses_stock (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    warehouse_uuid uuid NOT NULL,
    product_uuid uuid NOT NULL,
    start_quantity integer DEFAULT 0 NOT NULL,
    reserve_quantity integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);


--
-- TOC entry 3416 (class 0 OID 35613)
-- Dependencies: 217
-- Data for Name: inbox; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.inbox (uuid, metadata, status, type, created_at, updated_at, reference, response, retry_count, action) FROM stdin;
\.


--
-- TOC entry 3417 (class 0 OID 35629)
-- Dependencies: 218
-- Data for Name: outbox; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.outbox (uuid, metadata, status, type, created_at, updated_at, reference, response, retry_count, action) FROM stdin;
\.


--
-- TOC entry 3414 (class 0 OID 34957)
-- Dependencies: 215
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.warehouses (uuid, name, code, "desc", shop_uuid, created_at, updated_at, created_by, updated_by, status) FROM stdin;
1533773f-d3f5-47d0-b483-45646f030476	Gudang Gramedia Denpasar 1	gudang-gramedia-denpasar-1	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris ultrices est in nulla finibus facilisis. Vivamus neque tellus, suscipit eget purus sed, facilisis convallis libero. Quisque id eleifend risus. Quisque iaculis dapibus magna id sodales. Ut lacinia ante sed enim gravida, scelerisque commodo velit vulputate. Aliquam malesuada lacus ut nisi mollis condimentum. Fusce ac euismod velit. Morbi sit amet ullamcorper turpis. Phasellus ut cursus lorem. Fusce vitae tempor elit, sed viverra dui. Donec mattis, nisi sed gravida tincidunt, dolor mi venenatis urna, id ullamcorper dui nibh eu nunc.\n\n	ee7ae731-2eed-42f6-8852-d84c8847f112	2025-07-13 18:46:24.040278+00	2025-07-13 18:46:24.040278+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b	active
\.


--
-- TOC entry 3415 (class 0 OID 34971)
-- Dependencies: 216
-- Data for Name: warehouses_stock; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.warehouses_stock (uuid, warehouse_uuid, product_uuid, start_quantity, reserve_quantity, created_at, updated_at, created_by, updated_by) FROM stdin;
5087df34-1baf-4a5f-bd9c-ba066e2baf5a	1533773f-d3f5-47d0-b483-45646f030476	865edd5c-b190-4899-8c1a-e62780b7948b	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 18:51:58.631164+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
9d31ae13-55b8-4a73-9846-4c21044d1336	1533773f-d3f5-47d0-b483-45646f030476	9046175d-aaf8-4348-9ac3-3de1db6189e2	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 18:51:58.631164+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
8e7806ab-1376-4287-8dad-3ed2c22fcf8a	1533773f-d3f5-47d0-b483-45646f030476	e365ccc8-6b3d-4411-a2cb-a648f129f678	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 18:51:58.631164+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
c92b2690-bca5-4c44-8fad-a96c864f37e9	1533773f-d3f5-47d0-b483-45646f030476	6dad167b-6f72-48d2-a0c9-7747c0bbdacd	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 18:51:58.631164+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
ecf72266-e0ad-4a8d-9783-1feba0bb0e8c	1533773f-d3f5-47d0-b483-45646f030476	fc20e0db-0b06-41b2-9949-d0b4f57dec0b	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 18:51:58.631164+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
171ad69f-ca3c-46e0-8ad6-d8fb6eecfa90	1533773f-d3f5-47d0-b483-45646f030476	004d63cb-25f8-4a79-9b5f-2ee2d44ac709	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 21:52:43.783794+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
aa383a3b-898c-429c-814a-19e554609630	1533773f-d3f5-47d0-b483-45646f030476	96a70723-dc22-4ecd-abbf-b9024b736b22	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 22:59:06.978611+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
fd61a26a-04d2-4b09-987c-fa03e34140d4	1533773f-d3f5-47d0-b483-45646f030476	8e1ce434-1e1a-4692-98ef-b1f88afe7d2e	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 22:59:07.679729+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
40874710-a16d-41b4-89b2-d8579a18f0fb	1533773f-d3f5-47d0-b483-45646f030476	13744495-caa4-4e93-bdb2-68b68c0742e2	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 23:49:03.589216+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
8ae66e24-de46-4057-bfe9-2e1f0944faaa	1533773f-d3f5-47d0-b483-45646f030476	50ca5859-f78c-4c8a-a311-a8da117c6ab8	20	0	2025-07-13 18:51:58.631164+00	2025-07-13 23:49:03.59215+00	baf64f79-0bc2-4333-b126-c36445e6851b	baf64f79-0bc2-4333-b126-c36445e6851b
\.


--
-- TOC entry 3262 (class 2606 OID 35723)
-- Name: inbox inbox_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inbox
    ADD CONSTRAINT inbox_pk PRIMARY KEY (uuid);


--
-- TOC entry 3265 (class 2606 OID 35725)
-- Name: inbox inbox_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inbox
    ADD CONSTRAINT inbox_unique UNIQUE (type, reference, action);


--
-- TOC entry 3267 (class 2606 OID 35644)
-- Name: outbox outbox_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.outbox
    ADD CONSTRAINT outbox_pk PRIMARY KEY (uuid);


--
-- TOC entry 3270 (class 2606 OID 35727)
-- Name: outbox outbox_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.outbox
    ADD CONSTRAINT outbox_unique UNIQUE (type, reference, action);


--
-- TOC entry 3256 (class 2606 OID 34968)
-- Name: warehouses warehouse_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouse_pk PRIMARY KEY (uuid);


--
-- TOC entry 3259 (class 2606 OID 34980)
-- Name: warehouses_stock warehouse_stock_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses_stock
    ADD CONSTRAINT warehouse_stock_pk PRIMARY KEY (uuid);


--
-- TOC entry 3263 (class 1259 OID 35628)
-- Name: inbox_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX inbox_status_idx ON public.inbox USING btree (status);


--
-- TOC entry 3268 (class 1259 OID 35647)
-- Name: outbox_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX outbox_status_idx ON public.outbox USING btree (status);


--
-- TOC entry 3254 (class 1259 OID 34969)
-- Name: warehouse_code_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX warehouse_code_idx ON public.warehouses USING btree (code);


--
-- TOC entry 3257 (class 1259 OID 34970)
-- Name: warehouse_shop_uuid_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX warehouse_shop_uuid_idx ON public.warehouses USING btree (shop_uuid);


--
-- TOC entry 3260 (class 1259 OID 34981)
-- Name: warehouse_stock_warehouse_uuid_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX warehouse_stock_warehouse_uuid_idx ON public.warehouses_stock USING btree (warehouse_uuid, product_uuid);


-- Completed on 2025-07-15 00:37:26 WITA

--
-- PostgreSQL database dump complete
--

