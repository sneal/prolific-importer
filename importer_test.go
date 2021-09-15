package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func fixturePath(filename string) string {
	path, _ := filepath.Abs(filepath.Join("fixtures", filename))
	return path
}

func TestImportCsv(t *testing.T) {
	g := NewWithT(t)

	trackerProjectID := 12345
	trackerAPIToken := "not-a-real-api-token"

	expectedStories := []Story{
		{
			ProjectID: 12345,
			Name:      "TKG prerequisites are installed locally",
			StoryType: "feature",
			Labels:    []Label{{Name: "Install TKG"}},
			Description: `** As a person setting up TKG**
**I want ** to have the necessary tools installed to deploy TKG from my machine 
**So that ** I can begin deploying TKG


**Notes: The tools include the TKG CLI, kubectl, and docker. 
A list of them and instructions can be found in the documentation https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.1/vmware-tanzu-kubernetes-grid-11/GUID-install-tkg-set-up-tkg.html**`,
		},
		{
			ProjectID: 12345,
			Name:      "Vsphere environment is prepared",
			StoryType: "feature",
			Labels:    []Label{{Name: "Install TKG"}, {Name: "vsphere"}},
			Description: `** As a VI admin**
**I want ** to prepare my Vsphere environment 
**So that ** I can deploy TKG there.


**Notes: Documentation and detailed instructions are available: https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.1/vmware-tanzu-kubernetes-grid-11/GUID-install-tkg-vsphere.html**`,
		},
		{
			ProjectID: 12345,
			Name:      "Additional Vsphere prep tasks for disconnected environments",
			StoryType: "feature",
			Labels:    []Label{{Name: "Install TKG"}, {Name: "vsphere"}},
			Description: `** As a VI admin**
**I want ** to set up a docker registry and provide TKG's necessary resources locally 
**So that ** I can deploy TKG without internet access


**Notes: The bulk of this task is making the necessary docker images available offline. See the documented instructions here: https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.1/vmware-tanzu-kubernetes-grid-11/GUID-install-tkg-airgapped-environments.html**`,
		},
	}

	findExpectedStory := func(name string) *Story {
		for _, s := range expectedStories {
			if s.Name == name {
				return &s
			}
		}
		return nil
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var story Story
		err := decoder.Decode(&story)
		g.Expect(err).ShouldNot(HaveOccurred())

		expectedStory := findExpectedStory(story.Name)
		g.Expect(expectedStory).NotTo(BeNil(), "Could not find the expected story '%s'", story.Name)
		g.Expect(story.Name).To(Equal(expectedStory.Name))
		g.Expect(story.Description).To(Equal(expectedStory.Description))
		g.Expect(story.ProjectID).To(Equal(expectedStory.ProjectID))
		g.Expect(story.StoryType).To(Equal(expectedStory.StoryType))
		g.Expect(story.Labels).To(BeEquivalentTo(expectedStory.Labels))

		encoder := json.NewEncoder(w)
		err = encoder.Encode(story)
		g.Expect(err).ShouldNot(HaveOccurred())
	}))
	defer svr.Close()

	csv := fixturePath("stories1.csv")

	trackerGateway := NewTrackerGateway(svr.URL, trackerAPIToken)
	csvImporter := NewCsvImporter(trackerGateway)

	err := csvImporter.ImportCsvFromFile(trackerProjectID, csv)
	g.Expect(err).ShouldNot(HaveOccurred())
}
