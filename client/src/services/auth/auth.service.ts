import { LoginPayload } from "../../interfaces/LoginPayload.interface"
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
}
