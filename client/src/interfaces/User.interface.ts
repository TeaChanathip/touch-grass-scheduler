import { UUID } from "crypto"

export enum UserRole {
    STUDENT = "student",
    TEACHER = "teacher",
    GUARDIAN = "guardian",
    ADMIN = "admin",
}

export enum UserGender {
    MALE = "male",
    FEMALE = "female",
    OTHER = "other",
    PNTS = "prefer_not_to_say",
}

export interface User {
    id: UUID
    role: UserRole
    first_name: string
    middle_name?: string
    last_name?: string
    phone: string
    gender: UserGender
    email: string
    avatar_url?: string
    school_num?: string
}
