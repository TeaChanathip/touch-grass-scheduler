CREATE TYPE gender AS ENUM('male', 'female', 'other', 'prefer_not_to_say');
CREATE TYPE role AS ENUM('student', 'teacher', 'guardian', 'admin');

-- Entities
CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "role" role NOT NULL,
    "first_name" VARCHAR(128) NOT NULL,
    "middle_name" VARCHAR (128) DEFAULT '',
    "last_name" VARCHAR(128) DEFAULT '',
    "phone" VARCHAR(15) NOT NULL, -- e164
    "gender" gender NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" VARCHAR(60) NOT NULL,
    "avatar_key" VARCHAR(128) DEFAULT NULL, -- AKA. object name in the Storage
    "school_num" VARCHAR(16) DEFAULT NULL,
    UNIQUE("email")
);

CREATE TABLE IF NOT EXISTS "schools" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(128) NOT NULL,
    "class_count" INTEGER DEFAULT 0,
    "building_num" VARCHAR(16) NOT NULL,
    "moo" SMALLINT DEFAULT NULL,
    "soi" VARCHAR(32) DEFAULT NULL,
    "road" VARCHAR(32) NOT NULL,
    "sub_district" VARCHAR(32) NOT NULL,
    "district" VARCHAR(32) NOT NULL,
    "province" VARCHAR(32) NOT NULL
);

CREATE TABLE IF NOT EXISTS "classes" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(128) NOT NULL,
    "student_count" INTEGER DEFAULT 0,
    "school_id" UUID NOT NULL REFERENCES "schools"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "books" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "title" VARCHAR(128) NOT NULL,
    "page_count" INTEGER NOT NULL,
    "excercise_count" INTEGER NOT NULL,
    "total_man_hours" NUMERIC(4, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS "homework" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(128) NOT NULL,
    "man_hours" NUMERIC(3, 2) NOT NULL,
    "description" VARCHAR(1024) DEFAULT NULL,
    "score" DOUBLE PRECISION DEFAULT NULL,
    "book_id" UUID DEFAULT NULL REFERENCES "books"("id") ON DELETE CASCADE
);  

-- Relationships
CREATE TYPE relationship_type as ENUM('mother', 'father', 'other');

CREATE TABLE IF NOT EXISTS "student_guardians" (
    "student_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "guardian_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "type" relationship_type NOT NULL,
    PRIMARY KEY ("student_id", "guardian_id")
);

CREATE TABLE IF NOT EXISTS "class_teachers" (
    "class_id" UUID NOT NULL REFERENCES "classes"("id") ON DELETE CASCADE,
    "teacher_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    PRIMARY KEY ("class_id", "teacher_id")
);

CREATE TABLE IF NOT EXISTS "class_students" (
    "class_id" UUID NOT NULL REFERENCES "classes"("id") ON DELETE CASCADE,
    "student_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    PRIMARY KEY ("class_id", "student_id")
);

CREATE TABLE IF NOT EXISTS "homework_teachers" (
    "homework_id" UUID NOT NULL REFERENCES "homework"("id") ON DELETE CASCADE,
    "teacher_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    PRIMARY KEY ("homework_id", "teacher_id")
);

CREATE TABLE IF NOT EXISTS "homework_students" (
    "homework_id" UUID NOT NULL REFERENCES "homework"("id") ON DELETE CASCADE,
    "student_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "score" DOUBLE PRECISION DEFAULT NULL,
    PRIMARY KEY ("homework_id", "student_id")
);

CREATE TABLE IF NOT EXISTS "assignments" (
    "teacher_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "class_id" UUID NOT NULL REFERENCES "classes"("id") ON DELETE CASCADE,
    "homework_id" UUID NOT NULL REFERENCES "homework"("id") ON DELETE CASCADE,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "assigned_at" TIMESTAMPTZ DEFAULT NULL,
    "due_at" TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY ("teacher_id", "class_id", "homework_id")
);

-- For validating the files that is uploaded to object storage 
CREATE TYPE upload_type as ENUM('avatar');

CREATE TABLE IF NOT EXISTS "pending_uploads" (
    "object_key" VARCHAR(128) NOT NULL,
    "user_id" UUID NOT NULL REFERENCES "users"("id"),
    "type" upload_type NOT NULL,
    "expire_at" TIMESTAMPTZ NOT NULL,
    PRIMARY KEY ("object_key", "user_id")
);
