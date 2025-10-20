import { UserGender, UserRole } from "./User.interface"

export interface RegisterPayload {
    role: UserRole
    first_name: string
    middle_name?: string
    last_name?: string
    phone: string
    gender: UserGender
    password: string
    school_num?: string
}
