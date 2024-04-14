infra:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		deploy-infra \
				--path "./infra"

capi:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		install-capi \
				--cluster-name="dscott-capi"