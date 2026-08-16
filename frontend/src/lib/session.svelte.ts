import { api, type User } from '$lib/api';

class Session {
	user = $state<User | null>(null);
	ready = $state(false);
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
		this.ready = false;

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
		} finally {
			if (version === this.version) {
				this.ready = true;
			}
		}
	}

	setUser(user: User) {
		this.version += 1;
		this.user = user;
		this.ready = true;
	}

	clear() {
		this.version += 1;
		this.user = null;
		this.ready = true;
	}
}

export const session = new Session();
