# samba-container Release Process

## Preparation

The smbmetrics project has a dedicated branch, called `release`, for release
versions. This is done to update files, in particular the Dockerfile, that
control dependencies and versioning. Tags are applied directly to this branch
and only this branch.


### Tagging

Prior to tagging, we must update the `release` branch to "contain" all the
latest changes from the `main` branch. We do this by merging `main` into
`release`.
Example:

```
git checkout main
git pull --ff-only
git checkout release
git pull --ff-only
git merge main
# resolve any conflicts
```

Now we need to "pin" the appropriate version of the samba-server container
dependency.  Edit `Dockerfile` and change the tag part of the "FROM" line with
the `quay.io/samba.org/samba-server` image repository to use the latest
released samba-server tag.

At this point, an optional but recommended step is to do a test build before
tagging.  Run `make image-build`.

If you are happy with the content of the `release` branch, tag it. Example:

```
git checkout release
git tag -a -m 'Release v0.5' v0.5
```

This creates an annotated tag. Release tags must be annotated tags.

### Build

Using the tagged `release` branch, the container images for release will be
built. It is very important to ensure that base images are up-to-date.
It is very important to ensure that you perform the next set of steps with
clean new builds and do not use cached images. To accomplish both tasks it
is recommended to purge your local container engine of cached images
(Example: `podman image rm --all`). You should have no images named like
`quay.io/samba.org` in your local cache.

Build the images from scratch. Example:
```
make image-build
```

For the image that was just built, apply a temporary pre-release tag
to it. Example:
```
podman tag quay.io/samba.org/samba-metrics:{latest,v0.5pre1}
```

Log into quay.io.  Push the images to quay.io using the temporary tag. Example:
```
podman push quay.io/samba.org/samba-metrics:{latest,v0.5pre1}
```

Wait for the security scan to complete. There shouldn't be any issues if you
properly updated the base images before building. If there are issues and you
are sure you used the newest base images, check the base images on quay.io and
make sure that the number of issues are identical. The security scan can take
some time, while it runs you may want to do other things.


## GitHub Release

When you are satisfied that the tagged version is suitable for release, you
can push the tag to the public repo:
```
git push --follow-tags
```

Draft a new set of release notes. Select the recently pushed tag. Start with
the auto-generated release notes from GitHub (activate the `Generate release
notes` button/link). Add an introductory section (see previous notes for an
example). Add a "Highlights" section if there are any notable features or fixes
in the release. The Highlights section can be skipped if the content of the
release is unremarkable (e.g. few changes occurred since the previous release).

Because this is a container based release we do not provide any build artifacts
on GitHub (beyond the sources automatically provided there). Instead we add
a Downloads section that notes the exact tags and digests that the images can
be found at on quay.io.

Use the following partial snippet as an example:
```
## Download

Images built for this release can be obtained from the quay.io image registry.

* By tag: quay.io/samba.org/samba-metrics:v0.5
* By digest: quay.io/samba.org/samba-metrics@sha256:09c867343af39b237230f94a734eacc8313f2330c7d934994522ced46b740715
```
... using the image that was pushed earlier

The tag is pretty obvious - it should match the image tag (minus any pre-release
marker). You can get the digest from the tag using the quay.io UI (do not use
any local digest hashes). Click on the SHA256 link and then copy the full
manifest hash using the UI widget that appears.

Perform a final round of reviews, as needed, for the release notes and then
publish the release.

Once the release notes are drafted and then either immediately before or after
publishing them, use the quay.io UI to copy each pre-release tag to the "latest"
tag and a final "vX.Y" tag. Delete the temporary pre-release tags using the
quay.io UI as they are no longer needed.
