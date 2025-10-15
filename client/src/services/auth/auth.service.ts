import { LoginPayload } from "../../interfaces/LoginPayload.interface"
import { RegisterPayload } from "../../interfaces/RegisterPayload.interface"
import { User } from "../../interfaces/User.interface"
import { ApiService } from "../api.service"

export class AuthService {
    constructor(private readonly apiService: ApiService) {}

    private url = `${process.env.NEXT_PUBLIC_API_URL}/auth`

    async login(payload: LoginPayload) {
        return await this.apiService.post<
            LoginPayload,
            { user: User; token: string }
        >(`${this.url}/login`, payload)
    }

    async register(payload: RegisterPayload) {
        return await this.apiService.post<
            RegisterPayload,
            { user: User; token: string }
        >(`${this.url}/register`, payload)
    }
}
