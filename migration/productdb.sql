--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Debian 16.6-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Homebrew)

-- Started on 2025-07-20 12:14:00 WITA

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
-- TOC entry 3362 (class 1262 OID 34860)
-- Name: productdb; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE productdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


\connect productdb

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
-- TOC entry 3363 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 34893)
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying DEFAULT ''::character varying NOT NULL,
    "desc" text DEFAULT ''::text NOT NULL,
    top_image_url character varying DEFAULT ''::character varying NOT NULL,
    price double precision DEFAULT 0.0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    code character varying DEFAULT ''::character varying NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);


--
-- TOC entry 3356 (class 0 OID 34893)
-- Dependencies: 215
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.products VALUES ('96a70723-dc22-4ecd-abbf-b9024b736b22', 'Buku 1', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787361+00', '2025-07-13 18:43:41.787361+00', 'buku-1', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('8e1ce434-1e1a-4692-98ef-b1f88afe7d2e', 'Buku 2', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-2', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('004d63cb-25f8-4a79-9b5f-2ee2d44ac709', 'Buku 3', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-3', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('13744495-caa4-4e93-bdb2-68b68c0742e2', 'Buku 4', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-4', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('50ca5859-f78c-4c8a-a311-a8da117c6ab8', 'Buku 5', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-5', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('865edd5c-b190-4899-8c1a-e62780b7948b', 'Buku 6', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-6', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('9046175d-aaf8-4348-9ac3-3de1db6189e2', 'Buku 7', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-7', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('e365ccc8-6b3d-4411-a2cb-a648f129f678', 'Buku 8', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-8', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('6dad167b-6f72-48d2-a0c9-7747c0bbdacd', 'Buku 9', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-9', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');
INSERT INTO public.products VALUES ('fc20e0db-0b06-41b2-9949-d0b4f57dec0b', 'Buku 10', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lobortis tempus turpis quis vestibulum. Vestibulum lorem diam, sollicitudin eget arcu a, dapibus accumsan quam. Pellentesque blandit justo eu odio cursus, non mollis dui maximus. Suspendisse nibh tellus, luctus in risus ut, scelerisque faucibus augue. Mauris eget pretium orci. Quisque velit justo, scelerisque vitae rutrum eget, aliquet nec urna. Nam bibendum nunc metus, vitae aliquet urna faucibus feugiat. Mauris ac leo blandit elit semper finibus in in nulla. Nunc fermentum vel est malesuada faucibus.', 'https://miro.medium.com/v2/resize:fit:4800/format:webp/1*H90ZEKu0LMMIp96R6jbH3w.png', 20000, '2025-07-13 18:43:41.787+00', '2025-07-13 18:43:41.787+00', 'buku-10', 'baf64f79-0bc2-4333-b126-c36445e6851b', 'baf64f79-0bc2-4333-b126-c36445e6851b');


--
-- TOC entry 3212 (class 2606 OID 34907)
-- Name: products products_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pk PRIMARY KEY (uuid);


--
-- TOC entry 3210 (class 1259 OID 34908)
-- Name: products_code_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX products_code_idx ON public.products USING btree (code);


-- Completed on 2025-07-20 12:14:01 WITA

--
-- PostgreSQL database dump complete
--

