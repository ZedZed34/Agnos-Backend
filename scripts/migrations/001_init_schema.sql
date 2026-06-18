BEGIN;

-- ============================================================================
-- Source struct: domain.Hospital
--   ID   uint   `gorm:"primaryKey" json:"id"`
--   Name string `gorm:"unique;not null" json:"name"`
-- ============================================================================
CREATE TABLE IF NOT EXISTS hospitals (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT      NOT NULL UNIQUE
);

-- ============================================================================
-- Source struct: domain.Staff
--   ID           uint      `gorm:"primaryKey" json:"id"`
--   Username     string    `gorm:"unique" json:"username"`
--   PasswordHash string    `json:"-"`
--   HospitalID   uint      `json:"hospital_id"`
--   Hospital     Hospital  `gorm:"foreignKey:HospitalID" json:"hospital"`
--   CreatedAt    time.Time `json:"created_at"`
-- ============================================================================
CREATE TABLE IF NOT EXISTS staff (
    id            BIGSERIAL    PRIMARY KEY,
    username      TEXT         NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    hospital_id   BIGINT       NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_staff_hospital
        FOREIGN KEY (hospital_id) REFERENCES hospitals(id) ON DELETE CASCADE
);

-- ============================================================================
-- Source struct: domain.Patient
--   ID           uint      `gorm:"primaryKey" json:"id"`
--   FirstNameTH  string    `json:"first_name_th"`
--   MiddleNameTH string    `json:"middle_name_th"`
--   LastNameTH   string    `json:"last_name_th"`
--   FirstNameEN  string    `json:"first_name_en"`
--   MiddleNameEN string    `json:"middle_name_en"`
--   LastNameEN   string    `json:"last_name_en"`
--   PatientHN    string    `json:"patient_hn"`
--   NationalID   string    `gorm:"index" json:"national_id"`
--   PassportID   string    `gorm:"index" json:"passport_id"`
--   PhoneNumber  string    `json:"phone_number"`
--   Email        string    `json:"email"`
--   Gender       string    `json:"gender"`
--   DateOfBirth  string    `json:"date_of_birth"`
--   HospitalID   uint      `gorm:"index" json:"hospital_id"`
--   Hospital     Hospital  `gorm:"foreignKey:HospitalID" json:"hospital"`
--   CreatedAt    time.Time `json:"created_at"`
--   UpdatedAt    time.Time `json:"updated_at"`
-- ============================================================================
CREATE TABLE IF NOT EXISTS patients (
    id             BIGSERIAL    PRIMARY KEY,
    first_name_th  TEXT         NOT NULL,
    middle_name_th TEXT         NOT NULL,
    last_name_th   TEXT         NOT NULL,
    first_name_en  TEXT         NOT NULL,
    middle_name_en TEXT         NOT NULL,
    last_name_en   TEXT         NOT NULL,
    patient_hn     TEXT         NOT NULL,
    national_id    TEXT         NOT NULL,
    passport_id    TEXT         NOT NULL,
    phone_number   TEXT         NOT NULL,
    email          TEXT         NOT NULL,
    gender         TEXT         NOT NULL,
    date_of_birth  TEXT         NOT NULL,
    hospital_id    BIGINT       NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL,
    CONSTRAINT fk_patient_hospital
        FOREIGN KEY (hospital_id) REFERENCES hospitals(id) ON DELETE CASCADE
);

-- ----------------------------------------------------------------------------
-- Indexes for foreign key columns (lookups from child to parent)
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_staff_hospital_id    ON staff(hospital_id);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients(hospital_id);

-- ----------------------------------------------------------------------------
-- Indexes for search / lookup columns
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_patients_national_id ON patients(national_id);
CREATE INDEX IF NOT EXISTS idx_patients_passport_id ON patients(passport_id);

COMMIT;
