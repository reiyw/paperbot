package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromArxiv(t *testing.T) {
	paper, _ := FromArxivUrl("https://arxiv.org/abs/1805.09547")
	paper2, _ := FromArxivUrl("https://arxiv.org/pdf/1805.09547.pdf")
	assert.Equal(t, paper, paper2)
	assert.Equal(t, "1805.09547", paper.Id)
	assert.Equal(t, "Interpretable and Compositional Relation Learning by Joint Training with an Autoencoder", paper.Title)
	assert.Equal(t, []string{"Ryo Takahashi", "Ran Tian", "Kentaro Inui"}, paper.Authors)
	assert.Equal(t, "", paper.Volume)
	assert.Equal(t, "", paper.Venue)
	assert.Equal(t, 2018, paper.Year)
	assert.Equal(t, "https://arxiv.org/pdf/1805.09547.pdf", paper.PdfUrl)
	assert.Equal(t, "https://www.arxiv-vanity.com/papers/1805.09547/", paper.HtmlUrl)
	assert.Equal(t, `Embedding models for entities and relations are extremely useful for recovering missing facts in a knowledge base. Intuitively, a relation can be modeled by a matrix mapping entity vectors. However, relations reside on low dimension sub-manifolds in the parameter space of arbitrary matrices---for one reason, composition of two relations $\boldsymbol{M}_1,\boldsymbol{M}_2$ may match a third $\boldsymbol{M}_3$ (e.g. composition of relations currency_of_country and country_of_film usually matches currency_of_film_budget), which imposes compositional constraints to be satisfied by the parameters (i.e. $\boldsymbol{M}_1\cdot \boldsymbol{M}_2\approx \boldsymbol{M}_3$). In this paper we investigate a dimension reduction technique by training relations jointly with an autoencoder, which is expected to better capture compositional constraints. We achieve state-of-the-art on Knowledge Base Completion tasks with strongly improved Mean Rank, and show that joint training with an autoencoder leads to interpretable sparse codings of relations, helps discovering compositional constraints and benefits from compositional training. Our source code is released at github.com/tianran/glimvec.`, paper.AbstText)
	assert.Equal(t, "https://arxiv.org/abs/1805.09547", paper.AbstUrl)
	assert.Equal(t, "", paper.BibText)
	assert.Equal(t, "", paper.BibUrl)
	assert.Equal(t, "Equal contribution from first two authors. Accepted for publication in the ACL 2018", paper.Comment)
}

func TestFromAclweb(t *testing.T) {
	paper, _ := FromAclweb("https://aclanthology.info/papers/P18-1200/p18-1200")
	paper2, _ := FromAclweb("http://aclweb.org/anthology/P18-1200")
	assert.Equal(t, paper, paper2)
	assert.Equal(t, "P18-1200", paper.Id)
	assert.Equal(t, "Interpretable and Compositional Relation Learning by Joint Training with an Autoencoder", paper.Title)
	assert.Equal(t, []string{"Ryo Takahashi", "Ran Tian", "Kentaro Inui"}, paper.Authors)
	assert.Equal(t, "Proceedings of the 56th Annual Meeting of the Association for Computational Linguistics (Volume 1: Long Papers)", paper.Volume)
	assert.Equal(t, "ACL", paper.Venue)
	assert.Equal(t, 2018, paper.Year)
	assert.Equal(t, "http://aclweb.org/anthology/P18-1200", paper.PdfUrl)
	assert.Equal(t, "", paper.HtmlUrl)
	assert.Equal(t, "", paper.AbstText)
	assert.Equal(t, "https://aclanthology.info/papers/P18-1200/p18-1200", paper.AbstUrl)
	assert.Equal(t, "", paper.BibText)
	assert.Equal(t, "http://aclweb.org/anthology/P18-1200.bib", paper.BibUrl)
	assert.Equal(t, "", paper.Comment)
}
