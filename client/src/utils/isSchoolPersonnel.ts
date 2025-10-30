import { UserRole } from "../interfaces/User.interface"

export const schoolPersonnelRoles: UserRole[] = [
    UserRole.STUDENT,
    UserRole.TEACHER,
]

export function isSchoolPersonnel(role: UserRole): boolean {
    return schoolPersonnelRoles.includes(role)
}
