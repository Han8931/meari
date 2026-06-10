package core

import (
	"strings"
	"testing"
)

func TestIssueCertificate(t *testing.T) {
	svc := newCourseVault(t)

	meta, err := svc.IssueCertificate("balanced-trees", Certificate{
		CourseTitle: "Balanced Trees", Level: "intermediate", Date: "2026-06-10",
		Topics: 3, FirstTry: 2, Attempts: 4, Flawless: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if meta.Path != "meari-course/Balanced Trees/certificate.md" {
		t.Fatalf("certificate path = %q", meta.Path)
	}

	n, err := svc.OpenNote(meta.Path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Certificate of Completion", "Balanced Trees", "intermediate",
		"2026-06-10", "Topics finished: **3**", "2 / 3", "the echo returns",
	} {
		if !strings.Contains(n.Body, want) {
			t.Errorf("certificate missing %q:\n%s", want, n.Body)
		}
	}
	if strings.Contains(n.Body, "Flawless") {
		t.Error("a non-flawless certificate must not claim flawless")
	}

	// The certificate is a plain note, not a topic/manifest — the course still
	// loads with the same topics.
	if _, err := svc.LoadCourse("balanced-trees"); err != nil {
		t.Fatalf("course should still load after a certificate is added: %v", err)
	}
}

func TestIssueCertificateFlawless(t *testing.T) {
	svc := newCourseVault(t)
	meta, err := svc.IssueCertificate("balanced-trees", Certificate{
		CourseTitle: "Balanced Trees", Date: "2026-06-10", Topics: 3, FirstTry: 3, Attempts: 3, Flawless: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, _ := svc.OpenNote(meta.Path)
	if !strings.Contains(n.Body, "Flawless") {
		t.Errorf("flawless certificate should say so:\n%s", n.Body)
	}
}

func TestIssueCertificateUnknownCourse(t *testing.T) {
	svc := newCourseVault(t)
	if _, err := svc.IssueCertificate("no-such-course", Certificate{CourseTitle: "X"}); err == nil {
		t.Fatal("issuing a certificate for an unknown course should fail")
	}
}
