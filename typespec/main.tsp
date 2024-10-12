import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./error.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "GlassCMS API",
})
namespace GlassCMSCore;

alias CreateItem = TypeSpec.Rest.Resource.ResourceCreate<Item, Error>;
alias GetItems = TypeSpec.Rest.Resource.ResourceRead<Item, Error>;
alias UpdateItem = TypeSpec.Rest.Resource.ResourceUpdate<Item, Error>;

interface Items {
    @doc("Creates a new item.")
    op create is CreateItem.create;

    @doc("Retrieve an item by its unique identifier.")
    op get is GetItems.get;

    @doc("Update an item by its unique identifier.")
    op update is UpdateItem.update;
}

@doc("Item represents an individual content item.")
@resource("items")
model Item {
    @key
    @visibility("create", "read")
    name: string;
    display_name: string;
    content: string;
    create_time: utcDateTime;
    update_time: utcDateTime;
    @visibility("read")
    delete_time: utcDateTime;
    properties: Record<unknown>;
    metadata: Record<unknown>;
}
