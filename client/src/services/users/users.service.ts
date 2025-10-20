import { User } from "../../interfaces/User.interface"
import { ApiService } from "../api.service"

export class UsersService {
    constructor(private readonly apiService: ApiService) {}

    private url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/users`

    async getMe(): Promise<{ user: User }> {
        return await this.apiService.get<{ user: User }>(`${this.url}/me`)
    }
}
