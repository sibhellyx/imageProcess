package actions

import (
	"testing"
)

func TestDownload(t *testing.T) {

	t.Run("succsesfuly download", func(t *testing.T) {
		_, err := DownloadImageWithResty(
			"https://i.pinimg.com/736x/92/9f/3f/929f3fdd55668a80ff8a72b79ed911bd.jpg",
			"my_image1",
		)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("not image", func(t *testing.T) {
		_, err := DownloadImageWithResty(
			"https://www.google.com/url?sa=t&source=web&rct=j&opi=89978449&url=https://ftis-conf.mephi.ru/content/public/uploads/files/shablon_doklad.docx&ved=2ahUKEwiouuuC2M2PAxXtJRAIHV30DLcQFnoECBoQAQ&usg=AOvVaw1Nf6csOs0lPGdY9R90LfWG",
			"../../downloads/text",
		)

		t.Logf("Err : %v\n", err)

		if err == nil {
			t.Error("Expected error got")
		}
	})

}
