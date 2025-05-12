# Setting Up Your Project for Scoop

This document explains the changes made to enable installation of your project via the [Scoop](https://scoop.sh/) package manager for Windows.

## Files Added

1. **git-contrib.json** - The Scoop manifest file that describes how to install your application
2. **README.md** - Documentation including Scoop installation instructions
3. **.github/workflows/release.yml** - GitHub Actions workflow to automate building and releasing

## Next Steps

To make your project installable via Scoop, follow these steps:

1. **Push your project to GitHub**:
   ```
   git add .
   git commit -m "Add Scoop packaging support"
   git push
   ```

2. **Create a release**:
   - Create and push a tag (this will trigger the GitHub Actions workflow):
     ```
     git tag v1.0.0
     git push origin v1.0.0
     ```
   - The workflow will automatically build your application and create a release with the necessary files

3. **Update the manifest hash**:
   - After the release is created, download the `git-contrib-windows-amd64.zip.sha256` file
   - Update the `hash` field in `git-contrib.json` with the SHA256 hash from this file
   - Commit and push the updated manifest

4. **Test the installation**:
   ```
   scoop install https://raw.githubusercontent.com/username/git-contrib/main/git-contrib.json
   ```
   (Replace "username" with your GitHub username)

5. **Optional: Submit to a Scoop bucket**:
   - For wider distribution, you can submit your manifest to an existing Scoop bucket like "extras"
   - Follow the contribution guidelines of the specific bucket you want to submit to

## Maintenance

When releasing new versions:

1. Update the version in your code
2. Create and push a new tag (e.g., `v1.0.1`)
3. The GitHub Actions workflow will create a new release
4. Update the hash in the manifest file

## Resources

- [Scoop Documentation](https://github.com/ScoopInstaller/Scoop/wiki)
- [Scoop App Manifests](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)
- [Scoop Buckets](https://github.com/ScoopInstaller/Scoop/wiki/Buckets)