import { UpdateUserPayload, User } from "../../interfaces/User.interface"
import { ApiService } from "../api.service"

export class UsersService {
    constructor(private readonly apiService: ApiService) {}

    private url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/users`

    async getMe(): Promise<{ user: User }> {
        return await this.apiService.get<{ user: User }>(`${this.url}/me`)
    }

    async updateUser(
        updateUserPayload: UpdateUserPayload
    ): Promise<{ user: User }> {
        return await this.apiService.put<UpdateUserPayload, { user: User }>(
            this.url,
            updateUserPayload
        )
    }

    async getUploadAvartarSignedURL(): Promise<{
        url: string
        object_name: string
        form_data: { [key: string]: string }
    }> {
        return await this.apiService.get<{
            url: string
            object_name: string
            form_data: { [key: string]: string }
        }>(`${this.url}/avartar-signed-url`)
    }
}
