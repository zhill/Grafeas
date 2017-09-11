/* 
 * Grafeas API
 *
 * An API to insert and retrieve annotations on cloud artifacts.
 *
 * OpenAPI spec version: 0.1
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

// Source describes the location of the source used for the build.
type Source struct {

	// If provided, get the source from this location in in Google Cloud Storage.
	StorageSource StorageSource `json:"storageSource,omitempty"`

	// If provided, get source from this location in a Cloud Repo.
	RepoSource RepoSource `json:"repoSource,omitempty"`

	// If provided, the input binary artifacts for the build came from this location.
	ArtifactStorageSource StorageSource `json:"artifactStorageSource,omitempty"`

	// If provided, the source code used for the build came from this location.
	SourceContext ExtendedSourceContext `json:"sourceContext,omitempty"`

	// If provided, some of the source code used for the build may be found in these locations, in the case where the source repository had multiple remotes or submodules. This list will not include the context specified in the source_context field.
	AdditionalSourceContexts []ExtendedSourceContext `json:"additionalSourceContexts,omitempty"`

	// Hash(es) of the build source, which can be used to verify that the original source integrity was maintained in the build.  The keys to this map are file paths used as build source and the values contain the hash values for those files.  If the build source came in a single package such as a gzipped tarfile (.tar.gz), the FileHash will be for the single path to that file.
	FileHashes map[string]FileHashes `json:"fileHashes,omitempty"`
}
