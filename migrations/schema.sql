--
-- PostgreSQL database dump
--

\restrict 0ld1NNNYfyihcnoQKRisTxaIsXPgkVhYcijNDBwwohgqOs4fKKyvrDEC1O6Rntp

-- Dumped from database version 15.15 (Debian 15.15-1.pgdg13+1)
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: ai_conversations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_conversations (
    id integer NOT NULL,
    session_id character varying(255) NOT NULL,
    user_id integer,
    started_at timestamp with time zone DEFAULT now(),
    ended_at timestamp with time zone,
    total_messages integer DEFAULT 0,
    resulted_in_purchase boolean DEFAULT false,
    total_tokens_used integer DEFAULT 0,
    total_cost numeric(10,6) DEFAULT 0.00,
    user_agent text,
    ip_address character varying(45),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.ai_conversations OWNER TO postgres;

--
-- Name: TABLE ai_conversations; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ai_conversations IS 'Tracks AI assistant chat sessions';


--
-- Name: COLUMN ai_conversations.session_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_conversations.session_id IS 'Unique identifier for the chat session';


--
-- Name: COLUMN ai_conversations.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_conversations.user_id IS 'User ID if authenticated, null for anonymous';


--
-- Name: COLUMN ai_conversations.total_tokens_used; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_conversations.total_tokens_used IS 'Total OpenAI tokens consumed';


--
-- Name: COLUMN ai_conversations.total_cost; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_conversations.total_cost IS 'Estimated cost in USD';


--
-- Name: ai_conversations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_conversations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_conversations_id_seq OWNER TO postgres;

--
-- Name: ai_conversations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_conversations_id_seq OWNED BY public.ai_conversations.id;


--
-- Name: ai_feedback; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_feedback (
    id integer NOT NULL,
    message_id integer NOT NULL,
    conversation_id integer NOT NULL,
    helpful boolean,
    rating integer,
    feedback_text text,
    feedback_type character varying(50),
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT ai_feedback_feedback_type_check CHECK (((feedback_type)::text = ANY ((ARRAY['helpful'::character varying, 'not_helpful'::character varying, 'incorrect'::character varying, 'inappropriate'::character varying, 'other'::character varying])::text[]))),
    CONSTRAINT ai_feedback_rating_check CHECK (((rating >= 1) AND (rating <= 5)))
);


ALTER TABLE public.ai_feedback OWNER TO postgres;

--
-- Name: TABLE ai_feedback; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ai_feedback IS 'User feedback on AI assistant responses';


--
-- Name: COLUMN ai_feedback.helpful; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_feedback.helpful IS 'Was this response helpful?';


--
-- Name: COLUMN ai_feedback.rating; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_feedback.rating IS 'User rating from 1-5 stars';


--
-- Name: COLUMN ai_feedback.feedback_text; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_feedback.feedback_text IS 'Optional user comment';


--
-- Name: ai_feedback_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_feedback_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_feedback_id_seq OWNER TO postgres;

--
-- Name: ai_feedback_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_feedback_id_seq OWNED BY public.ai_feedback.id;


--
-- Name: ai_messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_messages (
    id integer NOT NULL,
    conversation_id integer NOT NULL,
    role character varying(20) NOT NULL,
    content text NOT NULL,
    tokens_used integer DEFAULT 0,
    response_time_ms integer,
    model character varying(50) DEFAULT 'gpt-3.5-turbo'::character varying,
    temperature numeric(3,2),
    metadata jsonb,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT ai_messages_role_check CHECK (((role)::text = ANY ((ARRAY['user'::character varying, 'assistant'::character varying, 'system'::character varying])::text[])))
);


ALTER TABLE public.ai_messages OWNER TO postgres;

--
-- Name: TABLE ai_messages; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ai_messages IS 'Stores individual messages in AI conversations';


--
-- Name: COLUMN ai_messages.role; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_messages.role IS 'Message sender: user, assistant, or system';


--
-- Name: COLUMN ai_messages.content; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_messages.content IS 'The actual message text';


--
-- Name: COLUMN ai_messages.tokens_used; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_messages.tokens_used IS 'OpenAI tokens consumed by this message';


--
-- Name: COLUMN ai_messages.response_time_ms; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_messages.response_time_ms IS 'API response time in milliseconds';


--
-- Name: COLUMN ai_messages.metadata; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_messages.metadata IS 'Additional data like product IDs mentioned, intents, etc.';


--
-- Name: ai_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_messages_id_seq OWNER TO postgres;

--
-- Name: ai_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_messages_id_seq OWNED BY public.ai_messages.id;


--
-- Name: ai_product_cache; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_product_cache (
    id integer NOT NULL,
    product_id integer NOT NULL,
    description_text text NOT NULL,
    search_keywords text[],
    category character varying(100),
    price_tier character varying(20),
    popularity_score integer DEFAULT 0,
    last_mentioned_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT ai_product_cache_price_tier_check CHECK (((price_tier)::text = ANY ((ARRAY['budget'::character varying, 'mid'::character varying, 'premium'::character varying])::text[])))
);


ALTER TABLE public.ai_product_cache OWNER TO postgres;

--
-- Name: TABLE ai_product_cache; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ai_product_cache IS 'Cached product information optimized for AI retrieval';


--
-- Name: COLUMN ai_product_cache.description_text; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_product_cache.description_text IS 'Product description formatted for AI context';


--
-- Name: COLUMN ai_product_cache.search_keywords; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_product_cache.search_keywords IS 'Keywords for semantic search';


--
-- Name: COLUMN ai_product_cache.popularity_score; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_product_cache.popularity_score IS 'How often product is mentioned in chats';


--
-- Name: ai_product_cache_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_product_cache_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_product_cache_id_seq OWNER TO postgres;

--
-- Name: ai_product_cache_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_product_cache_id_seq OWNED BY public.ai_product_cache.id;


--
-- Name: ai_user_preferences; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_user_preferences (
    id integer NOT NULL,
    user_id integer,
    session_id character varying(255),
    preferred_categories text[],
    budget_min numeric(10,2),
    budget_max numeric(10,2),
    interaction_count integer DEFAULT 0,
    last_products_viewed integer[],
    last_products_purchased integer[],
    conversation_style character varying(50),
    preferred_language character varying(10) DEFAULT 'en'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT ai_user_preferences_conversation_style_check CHECK (((conversation_style)::text = ANY ((ARRAY['concise'::character varying, 'detailed'::character varying, 'friendly'::character varying, 'professional'::character varying])::text[])))
);


ALTER TABLE public.ai_user_preferences OWNER TO postgres;

--
-- Name: TABLE ai_user_preferences; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ai_user_preferences IS 'Stores learned preferences from user interactions with AI';


--
-- Name: COLUMN ai_user_preferences.preferred_categories; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_user_preferences.preferred_categories IS 'Array of category names user is interested in';


--
-- Name: COLUMN ai_user_preferences.budget_min; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_user_preferences.budget_min IS 'Minimum budget mentioned by user';


--
-- Name: COLUMN ai_user_preferences.budget_max; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_user_preferences.budget_max IS 'Maximum budget mentioned by user';


--
-- Name: COLUMN ai_user_preferences.last_products_viewed; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_user_preferences.last_products_viewed IS 'Array of product IDs viewed in chat';


--
-- Name: COLUMN ai_user_preferences.conversation_style; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ai_user_preferences.conversation_style IS 'How the user prefers to communicate';


--
-- Name: ai_user_preferences_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_user_preferences_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_user_preferences_id_seq OWNER TO postgres;

--
-- Name: ai_user_preferences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_user_preferences_id_seq OWNED BY public.ai_user_preferences.id;


--
-- Name: customers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customers (
    id integer NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.customers OWNER TO postgres;

--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.customers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.customers_id_seq OWNER TO postgres;

--
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    widget_id integer,
    transaction_id integer,
    status_id integer,
    quantity integer NOT NULL,
    amount integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    customer_id integer
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    token character(43) NOT NULL,
    data bytea NOT NULL,
    expiry timestamp(6) without time zone NOT NULL
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: statuses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.statuses (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.statuses OWNER TO postgres;

--
-- Name: statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.statuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.statuses_id_seq OWNER TO postgres;

--
-- Name: statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.statuses_id_seq OWNED BY public.statuses.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tokens (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    expiry timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.tokens OWNER TO postgres;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tokens_id_seq OWNER TO postgres;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- Name: transaction_statuses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_statuses (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.transaction_statuses OWNER TO postgres;

--
-- Name: transaction_statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transaction_statuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transaction_statuses_id_seq OWNER TO postgres;

--
-- Name: transaction_statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transaction_statuses_id_seq OWNED BY public.transaction_statuses.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    amount integer NOT NULL,
    currency character varying(255) NOT NULL,
    last_four character varying(4) NOT NULL,
    bank_return_code character varying(255),
    transaction_status_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    expiry_month integer DEFAULT 0 NOT NULL,
    expiry_year integer DEFAULT 0 NOT NULL,
    payment_intent character varying(255) DEFAULT ''::character varying NOT NULL,
    payment_method character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_id_seq OWNER TO postgres;

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: widgets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.widgets (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    inventory_level integer,
    price integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    image character varying(255) DEFAULT ''::character varying NOT NULL,
    is_recurring boolean DEFAULT false NOT NULL,
    plan_id character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.widgets OWNER TO postgres;

--
-- Name: widgets_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.widgets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.widgets_id_seq OWNER TO postgres;

--
-- Name: widgets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.widgets_id_seq OWNED BY public.widgets.id;


--
-- Name: ai_conversations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_conversations ALTER COLUMN id SET DEFAULT nextval('public.ai_conversations_id_seq'::regclass);


--
-- Name: ai_feedback id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_feedback ALTER COLUMN id SET DEFAULT nextval('public.ai_feedback_id_seq'::regclass);


--
-- Name: ai_messages id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_messages ALTER COLUMN id SET DEFAULT nextval('public.ai_messages_id_seq'::regclass);


--
-- Name: ai_product_cache id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_product_cache ALTER COLUMN id SET DEFAULT nextval('public.ai_product_cache_id_seq'::regclass);


--
-- Name: ai_user_preferences id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_user_preferences ALTER COLUMN id SET DEFAULT nextval('public.ai_user_preferences_id_seq'::regclass);


--
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: statuses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.statuses ALTER COLUMN id SET DEFAULT nextval('public.statuses_id_seq'::regclass);


--
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- Name: transaction_statuses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_statuses ALTER COLUMN id SET DEFAULT nextval('public.transaction_statuses_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: widgets id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.widgets ALTER COLUMN id SET DEFAULT nextval('public.widgets_id_seq'::regclass);


--
-- Name: ai_conversations ai_conversations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_conversations
    ADD CONSTRAINT ai_conversations_pkey PRIMARY KEY (id);


--
-- Name: ai_conversations ai_conversations_session_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_conversations
    ADD CONSTRAINT ai_conversations_session_id_key UNIQUE (session_id);


--
-- Name: ai_feedback ai_feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_feedback
    ADD CONSTRAINT ai_feedback_pkey PRIMARY KEY (id);


--
-- Name: ai_messages ai_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT ai_messages_pkey PRIMARY KEY (id);


--
-- Name: ai_product_cache ai_product_cache_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_product_cache
    ADD CONSTRAINT ai_product_cache_pkey PRIMARY KEY (id);


--
-- Name: ai_product_cache ai_product_cache_product_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_product_cache
    ADD CONSTRAINT ai_product_cache_product_id_key UNIQUE (product_id);


--
-- Name: ai_user_preferences ai_user_preferences_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_user_preferences
    ADD CONSTRAINT ai_user_preferences_pkey PRIMARY KEY (id);


--
-- Name: ai_user_preferences ai_user_preferences_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_user_preferences
    ADD CONSTRAINT ai_user_preferences_user_id_key UNIQUE (user_id);


--
-- Name: customers customers_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_email_key UNIQUE (email);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (token);


--
-- Name: statuses statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.statuses
    ADD CONSTRAINT statuses_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_email_key UNIQUE (email);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_token_hash_key UNIQUE (token_hash);


--
-- Name: transaction_statuses transaction_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_statuses
    ADD CONSTRAINT transaction_statuses_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: widgets widgets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.widgets
    ADD CONSTRAINT widgets_pkey PRIMARY KEY (id);


--
-- Name: idx_ai_conversations_resulted_in_purchase; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_conversations_resulted_in_purchase ON public.ai_conversations USING btree (resulted_in_purchase);


--
-- Name: idx_ai_conversations_session_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_conversations_session_id ON public.ai_conversations USING btree (session_id);


--
-- Name: idx_ai_conversations_started_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_conversations_started_at ON public.ai_conversations USING btree (started_at);


--
-- Name: idx_ai_conversations_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_conversations_user_id ON public.ai_conversations USING btree (user_id);


--
-- Name: idx_ai_feedback_conversation_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_feedback_conversation_id ON public.ai_feedback USING btree (conversation_id);


--
-- Name: idx_ai_feedback_helpful; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_feedback_helpful ON public.ai_feedback USING btree (helpful);


--
-- Name: idx_ai_feedback_message_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_feedback_message_id ON public.ai_feedback USING btree (message_id);


--
-- Name: idx_ai_feedback_rating; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_feedback_rating ON public.ai_feedback USING btree (rating);


--
-- Name: idx_ai_messages_conversation_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_messages_conversation_id ON public.ai_messages USING btree (conversation_id);


--
-- Name: idx_ai_messages_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_messages_created_at ON public.ai_messages USING btree (created_at);


--
-- Name: idx_ai_messages_metadata; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_messages_metadata ON public.ai_messages USING gin (metadata);


--
-- Name: idx_ai_messages_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_messages_role ON public.ai_messages USING btree (role);


--
-- Name: idx_ai_product_cache_popularity_score; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_product_cache_popularity_score ON public.ai_product_cache USING btree (popularity_score DESC);


--
-- Name: idx_ai_product_cache_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_product_cache_product_id ON public.ai_product_cache USING btree (product_id);


--
-- Name: idx_ai_product_cache_search_keywords; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_product_cache_search_keywords ON public.ai_product_cache USING gin (search_keywords);


--
-- Name: idx_ai_user_preferences_session_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_user_preferences_session_id ON public.ai_user_preferences USING btree (session_id);


--
-- Name: idx_ai_user_preferences_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_user_preferences_user_id ON public.ai_user_preferences USING btree (user_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: sessions_expiry_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX sessions_expiry_idx ON public.sessions USING btree (expiry);


--
-- Name: ai_conversations ai_conversations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_conversations
    ADD CONSTRAINT ai_conversations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: ai_feedback ai_feedback_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_feedback
    ADD CONSTRAINT ai_feedback_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.ai_conversations(id) ON DELETE CASCADE;


--
-- Name: ai_feedback ai_feedback_message_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_feedback
    ADD CONSTRAINT ai_feedback_message_id_fkey FOREIGN KEY (message_id) REFERENCES public.ai_messages(id) ON DELETE CASCADE;


--
-- Name: ai_messages ai_messages_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT ai_messages_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.ai_conversations(id) ON DELETE CASCADE;


--
-- Name: ai_product_cache ai_product_cache_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_product_cache
    ADD CONSTRAINT ai_product_cache_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.widgets(id) ON DELETE CASCADE;


--
-- Name: ai_user_preferences ai_user_preferences_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_user_preferences
    ADD CONSTRAINT ai_user_preferences_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: orders fk_customer_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders fk_status_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders orders_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.transaction_statuses(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders orders_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders orders_widget_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_widget_id_fkey FOREIGN KEY (widget_id) REFERENCES public.widgets(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: transactions transactions_transaction_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_transaction_status_id_fkey FOREIGN KEY (transaction_status_id) REFERENCES public.transaction_statuses(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict 0ld1NNNYfyihcnoQKRisTxaIsXPgkVhYcijNDBwwohgqOs4fKKyvrDEC1O6Rntp

