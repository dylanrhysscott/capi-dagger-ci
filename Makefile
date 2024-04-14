deploy-plan:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		deploy-infra \
				--path "./infra" \
				--apply=false \
				stdout

deploy-infra:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		deploy-infra \
				--path "./infra" \
				stdout

install-capi:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		install-capi \
				--cluster-name="dscott-capi" \
				stdout

destroy-plan:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		deploy-infra \
				--path "./infra" \
				--destroy \
				--apply=false \
				stdout

destroy-infra:
	dagger call \
		with-docreds \
				--token=env:DIGITALOCEAN_ACCESS_TOKEN \
				--spaces-access-key=env:AWS_ACCESS_KEY_ID \
				--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY \
		deploy-infra \
				--path "./infra" \
				--destroy \
				stdout