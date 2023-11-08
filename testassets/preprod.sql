--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Debian 14.5-1.pgdg110+1)

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
-- Name: tiger; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA tiger;


ALTER SCHEMA tiger OWNER TO postgres;

--
-- Name: tiger_data; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA tiger_data;


ALTER SCHEMA tiger_data OWNER TO postgres;

--
-- Name: topology; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA topology;


ALTER SCHEMA topology OWNER TO postgres;

--
-- Name: SCHEMA topology; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA topology IS 'PostGIS Topology schema';


--
-- Name: fuzzystrmatch; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch WITH SCHEMA public;


--
-- Name: EXTENSION fuzzystrmatch; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION fuzzystrmatch IS 'determine similarities and distance between strings';


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


--
-- Name: postgis_tiger_geocoder; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder WITH SCHEMA tiger;


--
-- Name: EXTENSION postgis_tiger_geocoder; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis_tiger_geocoder IS 'PostGIS tiger geocoder and reverse geocoder';


--
-- Name: postgis_topology; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis_topology WITH SCHEMA topology;


--
-- Name: EXTENSION postgis_topology; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis_topology IS 'PostGIS topology spatial types and functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.accounts (
    id integer NOT NULL,
    name character varying(20) NOT NULL,
    password character varying(60) NOT NULL,
    admin boolean DEFAULT false NOT NULL,
    created_on timestamp without time zone NOT NULL
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.accounts_id_seq OWNER TO postgres;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: geopoints; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.geopoints (
    id integer NOT NULL,
    title character varying(30) NOT NULL,
    user_id integer NOT NULL,
    location public.geography(Point,4326),
    created_on timestamp without time zone NOT NULL,
    amplitudes double precision[],
    picture character varying(42) NOT NULL,
    sound character varying(42) NOT NULL,
    available boolean DEFAULT false NOT NULL
);


ALTER TABLE public.geopoints OWNER TO postgres;

--
-- Name: geopoints_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.geopoints_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.geopoints_id_seq OWNER TO postgres;

--
-- Name: geopoints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.geopoints_id_seq OWNED BY public.geopoints.id;


--
-- Name: geopoints_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.geopoints_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.geopoints_user_id_seq OWNER TO postgres;

--
-- Name: geopoints_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.geopoints_user_id_seq OWNED BY public.geopoints.user_id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: geopoints id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.geopoints ALTER COLUMN id SET DEFAULT nextval('public.geopoints_id_seq'::regclass);


--
-- Name: geopoints user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.geopoints ALTER COLUMN user_id SET DEFAULT nextval('public.geopoints_user_id_seq'::regclass);


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.accounts (id, name, password, admin, created_on) FROM stdin;
1	Preprod	$2a$10$pvsAF4W4kIbfMbXkEX7spung2QR4SyIZUHfbhr6nl8ltqmMc9gPeK	f	2022-12-25 09:21:39.17702
2	admin	$2a$10$OVzEd790/f7eGIwR.Rdm4ebiWKq2mDS5PVD0IyN5I4hD2vhRVgAam	t	2022-12-25 09:36:40.917079
\.


--
-- Data for Name: geopoints; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.geopoints (id, title, user_id, location, created_on, amplitudes, picture, sound, available) FROM stdin;
1	WithTemplate	1	0101000020E6100000969526A5A0DBF8BFB936548CF39B4740	2022-12-25 09:25:55.551	{0,3447,3949,-3189,-3192,2168,-3103,2357,-1567,-1399,-2688,-2168,-1664,3077,2605,2415,-1662,1621,-4355,5247,-4975,2166,1881,-1151,2371,1394,-3428,3633,-2389,-3199,2832,2943,-1893,-1876,-3449,-3158,-5211,-1913,-3455,2306,1397,5737,5415,-6232,6004,3344,4479,-1849,-3654,-2176,-6151,4473,3439,-3200,6412,-7022,-3387,6266,-11355,-6258,-11119,-7039,-4881,7551,3711,-3199,-1915,-4187,-2415,3688,3449,-2590,3710,3588,4221,-2338,-1905,-2844,2093,-2432,5991,5246,-3595,-3711,-1875,-3145,4172,-2943,-1147,-2124,-3446,2878,-3200,-2636,-3082,1865,-2679,1405,-1913,2159,1151,-2886,7039,-2882,-5753,3198,-1151,2130,-2176,-2353,-1662,1662,2089,-3585,1405,-2431,3915,3118,-1837,1660,-2174,-3450,2572,-2665,-1920,1314,-2939,-1403,-2686,1407,-2432,-2176,1383,1910,-1302,-3454,-1408,-3354,-2176,-2175,-2391,1407,5142,1919,-2688,2687,-1408,2150,-1909,-1914,2071,2101,-1152,-2431,1395,-1828,-4137,-1317,2416,-2341,2942,-3966,4697,2943,-3683,1907,-2912,1407,895,-1152,2169,1147,1105,-2820,-1152,3111,2686,1919,-2881,1406,2089,800,-1059,-2355,-1285,3606,-2678,1919,-1032,1662,-2061,2061,-1599,-1886,1551,-1407,-2175,1407,-1288,-1829,1149,-1045,-2674,1380,2684,-2337,2381,-2585,1904,-2176,1151,1151,2363,1407,-1613,1407,-2430,-1919,-3198,1663,1396,-1152,1645,2426,1663,1912,2333,-1664,1916,-1152,-896,1406,-896,1880,-2856,-3172,-1152,-2077,1658,-1919,-2432,1663,1344,-2346,2160,3192,-1152,-1408,-1657,1068,-1081,1814,-3406,2115,1908,3935,-3418,-1664,-4964,2098,5503,9029,-4479,-4736,1873,-2308,1146,3152,-1659,3425,2431,-1845,2922,-1570,1317,1306,-1664,-3145,-1911,-2420,-1919,-1603,2165,-1598,1301,-1912,628,-2432,-1041,-1603,895,-1901,1549,2072,-2174,1282,1660,823,2174,1655,813,1399,1663,1394,-2148,-1152,-1150,-2582,-2860,1132,1663,1638,1918,-1289,-3187,-1382,-1080,-1634,-1664,-1128,1407,1616,-2174,-1846,1405,2171,-1289,-2132,-1656,-3454,-1920,-1644,-2176,1407,-1845,-1664,1406,1605,1148,-1133,-1920,-1402,1663,2146,-2667,1151,1403,-3454,2171,1611,1579,-1104,-1539,-2429,-1915,-1301,2425,1132,-891,-1151,1150,-1644,-1296,-1596,1407,1853,-1663,-1081,-1644,-1664,-1349,2673,-2432,2678,-860,1061,1073,-1560,-2305,1377,-1553,-1152,-1664,-1150,-1911,-1322,-1664,1838,-1152,1105,-1659,-1085,-1918,-1659,-1894,-2332,-2352,-3188,-1152,-1851,-2341,-1405,-2103,-1398,2623,2686,2419,-1896,-2431,1855,895,-3704,-1919,-1577,-1408,-1663,-1909,1660,1577,1281,-1551,1663,-1604,-2176,1549,1919,2677,2126,-1151,-1662,-2160,1545,1574,-1405,-2176,-1316,-1142,639,-1664,-1919,1793,-2430,-2944,-1070,-3682,-2368,3195,-1109,-640,-1881,-1561,-1147,2678,1406,1883,-2432,-1920,2390,1560,-1408,-2944,2429,-1408,2325,3198,2942,-2431,-1399,-893,-2321,-894,1151,1112,1898,1919,895,1149,1042,-640,-1664,1294,870,-1407,-636,-1920,1389,-1408,-1408,-636,1347,-1663,-1858,1406,2101,-2097,-1407,2094,-1920,-1664,1151,1096,1407,-1366,-2631,-1836,2172,1832,1151,1126,2104,1067,1627,-1859,-1663,-1794,-1152,-1050,-1920,1566,-1398,-1326,-1135,-1874,-2318,-2170,-1408,-1407,-1405,1145,1342,-1407,-896,1915,-1151,-1149,-1920,1862,1404,-2345,-2430,-1658,-2176,-1150,2425,-1794,-2174,-1663,-1408,2170,1919,1846,-1625,-1815,-2394,-1848,1407,4143,1151,1663,1919,2675,2174,2418,1910,-2601,-1408,-1920,-1406,-3185,-2432,1663,1862,-1407,-2859,2685,-3654,-1847,-1407,1917,1285,1658,1293,-1919,2172,1654,1407,-1659,1551,1663,-2176,-1664,-1613,-2138,-1385,875,-1062,-1649,1613,-1282,2372,-1407,-1402,1918,1126,-896,1655,2174,-1919,2096,-1842,1151,-2611,1843,895,-896,-1119,-1663,-1542,-1152,1656,1352,-3199,1794,-1408,-1664,-1372,-2169,-2139,-1150,1143,-1330,2072,1663,-1904,1113,636,-1664,2430,-1916,-2580,1149,-2431,-1097,1151,-1658,1906,1071,1094,1919,-1404,-1152,2687,895,1662,1405,1883,2431}	clearing.webp	merle.aac	t
2	WithPicture	1	0101000020E610000018C75BC149ED0440787FBC57AD564640	2022-12-25 09:29:14.004	{0,1103,-3663,-2165,-1919,-1552,1855,-1408,-1663,-1664,1893,1057,866,1306,1885,-1404,791,1151,-3129,-1359,-1384,1620,1606,2687,1403,1795,1661,-1150,-2176,3455,1919,-1328,1919,2175,1090,1599,-1292,1594,1919,-1407,1653,1397,1406,2942,-3453,1663,-1628,-1825,-1657,1383,1151,2672,1914,-1918,-2432,-1636,-1151,-2400,1133,-2077,-1662,1130,-1635,1662,-1390,875,1132,-3167,-1889,-4224,1378,-2692,1663,-1810,-1919,1407,-2168,1149,1559,1826,-2175,2609,1861,-1663,-1408,-2432,-1894,-2432,-2430,-1408,3193,2856,2319,2175,5229,2419,2671,-2171,2419,1037,-1290,1663,1869,1660,-2106,-2147,1918,1918,-2160,-2172,2431,7745,-4446,-1660,1651,-1662,-2430,2685,-1664,2866,1105,1643,-1120,1582,-2432,1350,3625,2646,3928,-1894,2390,-2085,-4734,-1386,-11315,-1587,-2605,-2144,2175,-2915,-2049,1907,-1863,1544,3965,-2320,1658,1587,1660,-1039,-1407,-1300,2667,1660,1333,1617,1919,-4326,894,-2111,-877,1624,1898,2100,-1918,1904,1614,-1597,-1565,1369,895,1651,-2664,-896,-1794,1865,-2430,-2830,-2431,2419,2649,-896,-799,-2324,1406,2169,-2424,2169,778,1915,-1664,893,-895,3362,1648,-1908,-1408,-1836,1072,1538,639,-895,895,1388,1147,-1920,-1919,1150,-1152,-2144,-1661,-1918,-2678,-1151,1103,1387,-2176,-1151,-1346,-1840,-2367,-3148,-896,-1656,-1664,-896,566,-1657,-1613,1400,-1663,637,-1368,-893,1147,-3454,-1920,886,-2125,1537,-2112,-1909,639,-890,-1598,-1919,881,-1908,-896,-1140,-1132,-847,-1663,-1567,-1140,-1152,-2172,-1292,-1328,-1804,-1107,-2176,-2156,-1408,-1408,1283,-640,-1549,1291,605,-1393,1574,-1141,1323,-1660,-1664,1121,-1403,-1408,1811,-890,-639,2426,1610,-896,-1152,1663,1150,1555,1398,-815,-1150,893,1656,-1320,1407,-896,3434,2431,-2688,-2681,-1151,1407,-1856,1848,-1378,-1342,-1119,-2432,-2174,-1829,-1856,895,-1408,-1637,-2310,-1141,-1913,-2942,-2622,1094,892,-895,-1117,1308,1919,819,-1838,1404,1332,-801,-1330,-2428,-1392,-896,-1402,1911,2123,-2665,-1664,2397,1407,1404,1290,2316,1915,2174,-1916,-896,2149,2636,2388,-1643,1627,1293,1915,-3967,-896,1148,1553,1601,1654,1663,-1596,-1345,-1145,-1405,-1152,1407,894,-639,1328,-896,-1635,-1664,1608,-2071,-1920,-1659,1043,-1580,-893,-1586,-1150,-1661,-2174,1820,1151,1892,-1405,1658,-1149,-1326,1663,-1649,814,-1804,639,1395,-640,1135,-2139,1147,1406,-895,-1594,-1841,1882,-1152,-3115,-2174,-810,1149,1663,2171,-2167,-1151,-1396,-1573,-1054,-1138,-1904,1407,1404,1378,1136,-883,1661,1151,-2176,1151,895,1636,1151,845,-1307,-1853,-896,1919,1148,-635,-1408,-1318,-3109,1052,-1664,1405,-1657,-896,-1408,2074,-3127,1406,-1664,-808,-1400,-1306,-621,1151,1662,1401,-1556,1028,-1334,-2430,-611,1407,-2432,-1643,-1088,-1407,1113,1656,1151,-896,638,895,-2681,1407,-1029,-811,-1095,1151,-857,-1152,2156,-1297,802,-1615,-1408,617,1407,-1151,-1151,-2418,1130,1915,1395,-1326,-1334,-1664,1294,-1314,1150,-896,-1920,-1381,-896,791,-1628,1365,-2599,-1598,-1026,-1300,1151,-1587,1657,-1641,2175,1053,-1408,-1657,-2345,1407,-2322,1147,-2400,895,1632,-1091,-1589,1044,1101,1285,-1025,-1152,895,1051,1329,-896,-1627,1406,1659,-894,1919,-1659,628,2392,1392,1298,1151,1610,2127,-1585,1109,838,-1408,-1903,3128,-2430,-1144,-870,2623,-1152,-2368,-1152,1914,-1385,-1093,1147,2173,-1843,-1148,-1407,1149,1293,-2432,1658,1151,2114,1385,-1829,-1663,-877,-2925,-1920,1832,-895,-1340,-1150,-1401,1139,1138,-2174,-383,-2432,2368,1661,-883,1549,1834,639,2687,-2943,-3848,-1850,-896,-1084,-2594,-1107,1109,-1308,-1406,-1663,-896,-1052,1648,-1650,-830,2120,-1805,-2054,2155,1629,1661,-2307,2368,-1400,-1142,-1920,-2109,-1151,1918,-1146,-2400,1919,-2688,-2432,891,-2089,-640,-1920,-2572,-2942,-3184,1917,895,1566,-1821,-2128,1149,1540,-1918,-1805,-1152,-600}	russie.webp	merle.aac	t
\.


--
-- Data for Name: spatial_ref_sys; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.spatial_ref_sys (srid, auth_name, auth_srid, srtext, proj4text) FROM stdin;
\.


--
-- Data for Name: geocode_settings; Type: TABLE DATA; Schema: tiger; Owner: postgres
--

COPY tiger.geocode_settings (name, setting, unit, category, short_desc) FROM stdin;
\.


--
-- Data for Name: pagc_gaz; Type: TABLE DATA; Schema: tiger; Owner: postgres
--

COPY tiger.pagc_gaz (id, seq, word, stdword, token, is_custom) FROM stdin;
\.


--
-- Data for Name: pagc_lex; Type: TABLE DATA; Schema: tiger; Owner: postgres
--

COPY tiger.pagc_lex (id, seq, word, stdword, token, is_custom) FROM stdin;
\.


--
-- Data for Name: pagc_rules; Type: TABLE DATA; Schema: tiger; Owner: postgres
--

COPY tiger.pagc_rules (id, rule, is_custom) FROM stdin;
\.


--
-- Data for Name: topology; Type: TABLE DATA; Schema: topology; Owner: postgres
--

COPY topology.topology (id, name, srid, "precision", hasz) FROM stdin;
\.


--
-- Data for Name: layer; Type: TABLE DATA; Schema: topology; Owner: postgres
--

COPY topology.layer (topology_id, layer_id, schema_name, table_name, feature_column, feature_type, level, child_id) FROM stdin;
\.


--
-- Name: accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.accounts_id_seq', 2, true);


--
-- Name: geopoints_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.geopoints_id_seq', 2, true);


--
-- Name: geopoints_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.geopoints_user_id_seq', 1, false);


--
-- Name: accounts accounts_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_name_key UNIQUE (name);


--
-- Name: accounts accounts_password_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_password_key UNIQUE (password);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: geopoints geopoints_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.geopoints
    ADD CONSTRAINT geopoints_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

