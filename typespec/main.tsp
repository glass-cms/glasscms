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
alias ListItems = TypeSpec.Rest.Resource.ResourceList<Item, Error>;

@friendlyName("ItemUpsert")
model UpsertItem is TypeSpec.Rest.Resource.ResourceCreateModel<Item>;

interface Items {
    create is CreateItem.create;
    get is GetItems.get;
    update is UpdateItem.update;

    @autoRoute
    @doc("Lists all items")
    @listsResource(Item)
    list(@query fields?: string): Item[] | Error;
    
    @route("/items")
    @doc("Creates or updates many items")
    @createsOrUpdatesResource(Item)
    upsert(@body items: UpsertItem[]): Item[] | Error;
}

@doc("Item represents an individual content item.")
@resource("items")
model Item {
    @key
    @visibility("read", "create")
    name: string;

    display_name: string;
    content: string;
    create_time: utcDateTime;
    update_time: utcDateTime;

    @visibility("read")
    delete_time: utcDateTime;

    properties: Record<unknown>;
    metadata: Record<unknown>;

    @doc("represents a hash value calculated from the item's content.")
    @visibility("read")
    hash: string;
}
