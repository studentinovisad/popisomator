import { api, type User } from '$lib/api';

class Session {
	user = $state<User | null>(null);
	private currentUserRequest: Promise<User> | null = null;
	private version = 0;

	getCurrentUser() {
		if (!this.currentUserRequest) {
			this.currentUserRequest = api.currentUser().finally(() => {
				this.currentUserRequest = null;
			});
		}

		return this.currentUserRequest;
	}

	async refresh() {
		const version = this.version;

		try {
			const user = await this.getCurrentUser();
			if (version === this.version) {
				this.user = user;
			}
			return user;
		} catch {
			if (version === this.version) {
				this.user = null;
			}
			return null;
		}
	}

	setUser(user: User) {
		this.version += 1;
		this.user = user;
	}

	clear() {
		this.version += 1;
		this.user = null;
	}
}

export const session = new Session();
